// SPDX-License-Identifier: MIT
pragma solidity 0.8.28;

import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/token/ERC721/IERC721.sol";
import "@openzeppelin/contracts/token/ERC721/IERC721Receiver.sol";
import "@openzeppelin/contracts/token/common/ERC2981.sol";
import "@openzeppelin/contracts/interfaces/IERC2981.sol";
import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/utils/Pausable.sol";
import "@openzeppelin/contracts/utils/ReentrancyGuard.sol";

/**
 * @title SudaMarketplace
 * @notice NFT marketplace with escrow-based listings.
 *
 * Flow:
 *   1. Seller calls `list(tokenId, price)`.
 *      The NFT is transferred from seller into the marketplace contract
 *      (escrow). The contract becomes the temporary holder.
 *   2. Buyer calls `buy(listingId)`.
 *      The contract:
 *        - pulls SUDA from buyer (ERC20.transferFrom),
 *        - splits the price into (royalty / market fee / seller share),
 *        - transfers NFT from itself to buyer.
 *      All three movements happen atomically.
 *   3. Seller can `cancel(listingId)` while ACTIVE — NFT returns to seller.
 *
 * Properties:
 *   - Uses ERC-2981 for royalty (asks the NFT contract via royaltyInfo).
 *     If the NFT doesn't support 2981, royalty = 0 (graceful degradation).
 *   - Market fee: 1% (configurable, hard-capped at 10%) goes to owner (treasury).
 *   - Pausable: emergency stop on listing/buying. Cancel is always allowed
 *     so sellers can withdraw NFTs even when paused.
 *   - ReentrancyGuard: protects buy() since it makes external calls.
 */
contract SudaMarketplace is IERC721Receiver, Ownable, Pausable, ReentrancyGuard {
    /// @notice Listing status.
    enum Status { ACTIVE, SOLD, CANCELLED }

    struct Listing {
        address seller;
        address nftContract;
        uint256 tokenId;
        uint256 price;
        Status  status;
    }

    /// @notice The SUDA ERC-20 token used as currency. Immutable.
    IERC20 public immutable sudaToken;

    /// @notice Marketplace fee in basis points (100 = 1%).
    uint256 public marketplaceFeeBps;

    /// @notice Hard cap on the market fee — owner cannot exceed this.
    uint256 public constant MAX_FEE_BPS = 1000; // 10%

    /// @notice listingId → Listing
    mapping(uint256 => Listing) public listings;
    uint256 public nextListingId;

    event Listed(
        uint256 indexed listingId,
        address indexed seller,
        address indexed nftContract,
        uint256 tokenId,
        uint256 price
    );
    event Sold(
        uint256 indexed listingId,
        address indexed buyer,
        address indexed seller,
        uint256 price,
        uint256 royaltyAmount,
        uint256 marketFee,
        uint256 sellerProceeds
    );
    event Cancelled(uint256 indexed listingId, address indexed seller);
    event MarketplaceFeeUpdated(uint256 oldBps, uint256 newBps);

    constructor(address initialOwner, IERC20 _sudaToken)
        Ownable(initialOwner)
    {
        require(address(_sudaToken) != address(0), "marketplace: zero token");
        sudaToken = _sudaToken;
        marketplaceFeeBps = 100; // 1% default
    }

    // ────────────────────────────────────────────────────────
    // Listing operations
    // ────────────────────────────────────────────────────────

    /**
     * @notice List an NFT for sale. The NFT is escrowed by the contract.
     *         Caller must `approve` the marketplace on the NFT contract first.
     * @return listingId Sequential listing ID (starts at 1).
     */
    function list(address nftContract, uint256 tokenId, uint256 price)
        external
        whenNotPaused
        returns (uint256 listingId)
    {
        require(price > 0, "marketplace: zero price");
        require(nftContract != address(0), "marketplace: zero nft");

        // Pull NFT into escrow. Will revert if not approved or not owned.
        IERC721(nftContract).safeTransferFrom(msg.sender, address(this), tokenId);

        unchecked {
            listingId = ++nextListingId;
        }
        listings[listingId] = Listing({
            seller:      msg.sender,
            nftContract: nftContract,
            tokenId:     tokenId,
            price:       price,
            status:      Status.ACTIVE
        });

        emit Listed(listingId, msg.sender, nftContract, tokenId, price);
    }

    /**
     * @notice Buy an active listing. Caller must `approve` SUDA spending first.
     */
    function buy(uint256 listingId) external whenNotPaused nonReentrant {
        Listing storage l = listings[listingId];
        require(l.status == Status.ACTIVE, "marketplace: not active");
        require(msg.sender != l.seller, "marketplace: cannot buy own listing");

        // Mark sold BEFORE external calls (CEI pattern).
        l.status = Status.SOLD;

        uint256 price = l.price;

        // 1) Calculate royalty via EIP-2981 (graceful if not supported).
        (address royaltyReceiver, uint256 royaltyAmount) =
            _getRoyalty(l.nftContract, l.tokenId, price);

        // 2) Calculate market fee.
        uint256 marketFee = (price * marketplaceFeeBps) / 10000;

        // 3) Whatever's left goes to the seller.
        uint256 sellerProceeds = price - royaltyAmount - marketFee;

        // 4) Pull total price from buyer once.
        require(
            sudaToken.transferFrom(msg.sender, address(this), price),
            "marketplace: SUDA transferFrom failed"
        );

        // 5) Distribute. Each transfer is from the contract's own balance.
        if (royaltyAmount > 0 && royaltyReceiver != address(0)) {
            require(
                sudaToken.transfer(royaltyReceiver, royaltyAmount),
                "marketplace: royalty payout failed"
            );
        }
        if (marketFee > 0) {
            require(
                sudaToken.transfer(owner(), marketFee),
                "marketplace: fee payout failed"
            );
        }
        require(
            sudaToken.transfer(l.seller, sellerProceeds),
            "marketplace: seller payout failed"
        );

        // 6) Send NFT to buyer.
        IERC721(l.nftContract).safeTransferFrom(address(this), msg.sender, l.tokenId);

        emit Sold(
            listingId,
            msg.sender,
            l.seller,
            price,
            royaltyAmount,
            marketFee,
            sellerProceeds
        );
    }

    /**
     * @notice Cancel an active listing. Only the seller can call.
     *         Allowed even when paused — sellers must always be able to
     *         retrieve their NFTs.
     */
    function cancel(uint256 listingId) external nonReentrant {
        Listing storage l = listings[listingId];
        require(l.status == Status.ACTIVE, "marketplace: not active");
        require(msg.sender == l.seller, "marketplace: not seller");

        l.status = Status.CANCELLED;

        IERC721(l.nftContract).safeTransferFrom(address(this), l.seller, l.tokenId);

        emit Cancelled(listingId, l.seller);
    }

    // ────────────────────────────────────────────────────────
    // Admin
    // ────────────────────────────────────────────────────────

    function setMarketplaceFeeBps(uint256 newBps) external onlyOwner {
        require(newBps <= MAX_FEE_BPS, "marketplace: fee too high");
        emit MarketplaceFeeUpdated(marketplaceFeeBps, newBps);
        marketplaceFeeBps = newBps;
    }

    function pause() external onlyOwner {
        _pause();
    }

    function unpause() external onlyOwner {
        _unpause();
    }

    // ────────────────────────────────────────────────────────
    // Internal helpers
    // ────────────────────────────────────────────────────────

    /**
     * @dev Returns (royaltyReceiver, royaltyAmount) for the given NFT,
     *      using EIP-2981. If the NFT contract doesn't support 2981,
     *      returns (zero, zero) so we silently skip royalty.
     */
    function _getRoyalty(address nftContract, uint256 tokenId, uint256 salePrice)
        internal
        view
        returns (address, uint256)
    {
        try IERC2981(nftContract).royaltyInfo(tokenId, salePrice)
            returns (address receiver, uint256 amount)
        {
            // Sanity: never let royalty exceed half the sale price.
            if (amount > salePrice / 2) {
                return (address(0), 0);
            }
            return (receiver, amount);
        } catch {
            return (address(0), 0);
        }
    }

    /// @inheritdoc IERC721Receiver
    function onERC721Received(address, address, uint256, bytes calldata)
        external
        pure
        override
        returns (bytes4)
    {
        return IERC721Receiver.onERC721Received.selector;
    }
}