// SPDX-License-Identifier: MIT
pragma solidity 0.8.28;

import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/utils/Pausable.sol";
import "@openzeppelin/contracts/utils/ReentrancyGuard.sol";

/**
 * @title SudaFundraising
 * @notice Crowdfunding-style fundraisers in SUDA.
 *
 *   Lifecycle:
 *     ACTIVE ──donate()──► ACTIVE  (raised += amount; status auto-flips to GOAL_REACHED when raised >= goal)
 *     GOAL_REACHED ──withdraw()──► WITHDRAWN  (creator only, sends raised to creator)
 *     ACTIVE ──(deadline passes && raised < goal)──► EXPIRED  (set on first refund)
 *     EXPIRED ──refund()──► EXPIRED  (each donor pulls their own contribution)
 *
 *   Key design points:
 *     - "Withdraw immediately on goal reached" — creator doesn't wait for deadline.
 *     - If donations exceed the goal, all of them go to the creator (no overflow refund).
 *     - Refund is pull-based: each donor calls refund() themselves. This is gas-safe
 *       and avoids the contract iterating over an unbounded list.
 *     - A donor whose refund hasn't been collected can still call refund any time;
 *       the contract just maps (fundraiser, donor) → contributed amount and zeroes
 *       it out on refund.
 */
contract SudaFundraising is Ownable, Pausable, ReentrancyGuard {
    enum Status {
        NONE,
        ACTIVE,
        GOAL_REACHED,
        WITHDRAWN,
        EXPIRED
    }

    struct Fundraiser {
        address creator;
        uint256 goal;
        uint256 raised;
        uint256 deadline;
        Status  status;
    }

    /// @notice The SUDA ERC-20 token used as currency.
    IERC20 public immutable sudaToken;

    /// @notice fundraiserId → Fundraiser
    mapping(uint256 => Fundraiser) public fundraisers;

    /// @notice fundraiserId → donor → amount contributed (cumulative).
    mapping(uint256 => mapping(address => uint256)) public contributions;

    uint256 public nextFundraiserId;

    event FundraiserCreated(
        uint256 indexed fundraiserId,
        address indexed creator,
        uint256 goal,
        uint256 deadline
    );
    event Donated(
        uint256 indexed fundraiserId,
        address indexed donor,
        uint256 amount,
        uint256 totalRaised
    );
    event GoalReached(uint256 indexed fundraiserId, uint256 totalRaised);
    event Withdrawn(uint256 indexed fundraiserId, address indexed creator, uint256 amount);
    event Expired(uint256 indexed fundraiserId);
    event Refunded(uint256 indexed fundraiserId, address indexed donor, uint256 amount);

    constructor(address initialOwner, IERC20 _sudaToken)
        Ownable(initialOwner)
    {
        require(address(_sudaToken) != address(0), "fundraising: zero token");
        sudaToken = _sudaToken;
    }

    // ────────────────────────────────────────────────────────
    // Lifecycle
    // ────────────────────────────────────────────────────────

    /// @notice Create a new fundraiser.
    function create(uint256 goal, uint256 deadline)
        external
        whenNotPaused
        returns (uint256 fundraiserId)
    {
        require(goal > 0, "fundraising: zero goal");
        require(deadline > block.timestamp, "fundraising: deadline in past");

        unchecked {
            fundraiserId = ++nextFundraiserId;
        }
        fundraisers[fundraiserId] = Fundraiser({
            creator:  msg.sender,
            goal:     goal,
            raised:   0,
            deadline: deadline,
            status:   Status.ACTIVE
        });

        emit FundraiserCreated(fundraiserId, msg.sender, goal, deadline);
    }

    /**
     * @notice Donate SUDA to a fundraiser. Caller must `approve` SUDA spending first.
     *         If the cumulative raised amount crosses the goal, status flips to GOAL_REACHED.
     */
    function donate(uint256 fundraiserId, uint256 amount)
        external
        whenNotPaused
        nonReentrant
    {
        Fundraiser storage f = fundraisers[fundraiserId];
        require(
            f.status == Status.ACTIVE || f.status == Status.GOAL_REACHED,
            "fundraising: not accepting donations"
        );
        require(block.timestamp < f.deadline, "fundraising: past deadline");
        require(amount > 0, "fundraising: zero amount");

        require(
            sudaToken.transferFrom(msg.sender, address(this), amount),
            "fundraising: SUDA transferFrom failed"
        );

        f.raised += amount;
        contributions[fundraiserId][msg.sender] += amount;

        emit Donated(fundraiserId, msg.sender, amount, f.raised);

        // Auto-flip status on first crossing of the goal.
        if (f.status == Status.ACTIVE && f.raised >= f.goal) {
            f.status = Status.GOAL_REACHED;
            emit GoalReached(fundraiserId, f.raised);
        }
    }

    /**
     * @notice Withdraw all raised funds. Only the creator. Only when GOAL_REACHED.
     *         Sends the entire `raised` amount (including any donations beyond the goal).
     */
    function withdraw(uint256 fundraiserId) external nonReentrant {
        Fundraiser storage f = fundraisers[fundraiserId];
        require(f.status == Status.GOAL_REACHED, "fundraising: not goal reached");
        require(msg.sender == f.creator, "fundraising: not creator");

        uint256 amount = f.raised;
        f.status = Status.WITHDRAWN;

        require(
            sudaToken.transfer(f.creator, amount),
            "fundraising: payout failed"
        );

        emit Withdrawn(fundraiserId, f.creator, amount);
    }

    /**
     * @notice Refund a donation from an expired fundraiser.
     *         Anyone can call to mark the fundraiser EXPIRED if the deadline
     *         has passed and the goal was not reached. Each donor calls
     *         refund() themselves to retrieve their contribution.
     */
    function refund(uint256 fundraiserId) external nonReentrant {
        Fundraiser storage f = fundraisers[fundraiserId];
        require(f.status == Status.ACTIVE || f.status == Status.EXPIRED,
            "fundraising: not refundable");
        require(block.timestamp >= f.deadline, "fundraising: not yet expired");
        require(f.raised < f.goal, "fundraising: goal reached, no refunds");

        // Mark expired on first refund call.
        if (f.status == Status.ACTIVE) {
            f.status = Status.EXPIRED;
            emit Expired(fundraiserId);
        }

        uint256 amount = contributions[fundraiserId][msg.sender];
        require(amount > 0, "fundraising: nothing to refund");

        contributions[fundraiserId][msg.sender] = 0;

        require(
            sudaToken.transfer(msg.sender, amount),
            "fundraising: refund failed"
        );

        emit Refunded(fundraiserId, msg.sender, amount);
    }

    // ────────────────────────────────────────────────────────
    // Admin
    // ────────────────────────────────────────────────────────

    function pause() external onlyOwner {
        _pause();
    }

    function unpause() external onlyOwner {
        _unpause();
    }
}