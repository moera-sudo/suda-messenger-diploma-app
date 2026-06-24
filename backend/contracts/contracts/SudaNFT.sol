// SPDX-License-Identifier: MIT
pragma solidity 0.8.28;

import "@openzeppelin/contracts/token/ERC721/ERC721.sol";
import "@openzeppelin/contracts/token/ERC721/extensions/ERC721URIStorage.sol";
import "@openzeppelin/contracts/token/ERC721/extensions/ERC721Pausable.sol";
import "@openzeppelin/contracts/token/common/ERC2981.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

/**
 * @title SudaNFT
 * @notice ERC-721 NFT collection for stickers, gifts, avatar frames.
 *
 * Properties:
 *   - One contract for the whole project (categories live off-chain in DB).
 *   - tokenURI is HTTP, points to Transaction Service API endpoint
 *     that returns OpenSea-style JSON metadata. The JSON is built from
 *     tx_nft_items + media-service.
 *   - Only the owner (treasury / backend) can mint. Users do not call mint
 *     directly — they request stickers/gifts via UI, the backend mints
 *     after validating the SUDA payment off-chain (in custodial v1) or
 *     on-chain via the marketplace contract (later).
 *   - EIP-2981 royalty: 5% on every secondary sale, paid to treasury.
 *     Marketplaces that respect EIP-2981 (including SudaMarketplace) will
 *     deduct this from the seller's proceeds.
 *   - Pausable: owner can freeze all transfers in emergencies.
 */
contract SudaNFT is ERC721, ERC721URIStorage, ERC721Pausable, ERC2981, Ownable {
    /// @notice Counter for sequential token IDs (1-based).
    uint256 private _nextTokenId;

    /// @notice Royalty in basis points (500 = 5.00%).
    uint96 public constant ROYALTY_BPS = 500;

    event Minted(address indexed to, uint256 indexed tokenId, string tokenURI);

    /**
     * @param initialOwner Address that owns the contract and receives royalties.
     */
    constructor(address initialOwner)
        ERC721("Suda NFT", "SUDA-NFT")
        Ownable(initialOwner)
    {
        // Default royalty for every token: 5% to treasury.
        _setDefaultRoyalty(initialOwner, ROYALTY_BPS);
    }

    /**
     * @notice Mint a new NFT to a user. Only owner (treasury) can call.
     * @param to       Recipient address.
     * @param uri      Metadata URI (HTTP to backend, e.g.
     *                 https://api.suda.app/api/v1/tx/nft/<id>/metadata).
     * @return tokenId Newly minted token ID (1, 2, 3, ...).
     */
    function mintTo(address to, string calldata uri)
        external
        onlyOwner
        returns (uint256 tokenId)
    {
        unchecked {
            tokenId = ++_nextTokenId;
        }
        _safeMint(to, tokenId);
        _setTokenURI(tokenId, uri);
        emit Minted(to, tokenId, uri);
    }

    /// @notice Total number of NFTs minted so far (NOT a supply cap).
    function totalMinted() external view returns (uint256) {
        return _nextTokenId;
    }

    function pause() external onlyOwner {
        _pause();
    }

    function unpause() external onlyOwner {
        _unpause();
    }

    // ────────────────────────────────────────────────────────
    // Required overrides (multiple inheritance)
    // ────────────────────────────────────────────────────────

    function tokenURI(uint256 tokenId)
        public
        view
        override(ERC721, ERC721URIStorage)
        returns (string memory)
    {
        return super.tokenURI(tokenId);
    }

    function supportsInterface(bytes4 interfaceId)
        public
        view
        override(ERC721, ERC721URIStorage, ERC2981)
        returns (bool)
    {
        return super.supportsInterface(interfaceId);
    }

    function _update(address to, uint256 tokenId, address auth)
        internal
        override(ERC721, ERC721Pausable)
        returns (address)
    {
        return super._update(to, tokenId, auth);
    }
}