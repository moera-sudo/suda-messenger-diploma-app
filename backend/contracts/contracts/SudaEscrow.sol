// SPDX-License-Identifier: MIT
pragma solidity 0.8.28;

import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/utils/Pausable.sol";
import "@openzeppelin/contracts/utils/ReentrancyGuard.sol";

/**
 * @title SudaEscrow
 * @notice Quest / bounty escrow.
 *
 *   Lifecycle:
 *     OPEN ──claim()──► CLAIMED ──submit()──► SUBMITTED ──approve()──► APPROVED
 *      │                  │
 *      └─cancel()────────┴────► CANCELLED   (creator only, before APPROVED)
 *      │
 *      └─expire()──────────────► EXPIRED    (anyone, after deadline, before APPROVED)
 *
 *   - createQuest pulls SUDA from creator into the contract (escrow).
 *   - approve sends SUDA to the assignee. Terminal.
 *   - cancel / expire send SUDA back to the creator. Terminal.
 *   - No dispute resolution in v1: if the creator refuses to approve,
 *     the deadline mechanism eventually returns funds to creator.
 *
 *   Trust model: this is the part of the system the user genuinely owns.
 *   The backend holds users' private keys (custodial), but it cannot pull
 *   funds out of escrow except via the rules encoded here. APPROVE requires
 *   the creator's signature, EXPIRE requires the deadline to have passed.
 */
contract SudaEscrow is Ownable, Pausable, ReentrancyGuard {
    enum Status {
        NONE,       // sentinel for non-existent quests
        OPEN,
        CLAIMED,
        SUBMITTED,
        APPROVED,
        CANCELLED,
        EXPIRED
    }

    struct Quest {
        address creator;
        address assignee;       // zero until claimed
        uint256 reward;         // SUDA amount in escrow
        uint256 deadline;       // unix timestamp
        Status  status;
    }

    /// @notice The SUDA ERC-20 token used as currency.
    IERC20 public immutable sudaToken;

    /// @notice questId → Quest
    mapping(uint256 => Quest) public quests;
    uint256 public nextQuestId;

    event QuestCreated(
        uint256 indexed questId,
        address indexed creator,
        uint256 reward,
        uint256 deadline
    );
    event QuestClaimed(uint256 indexed questId, address indexed assignee);
    event QuestSubmitted(uint256 indexed questId, address indexed assignee);
    event QuestApproved(uint256 indexed questId, address indexed assignee, uint256 reward);
    event QuestCancelled(uint256 indexed questId, address indexed creator, uint256 refund);
    event QuestExpired(uint256 indexed questId, uint256 refund);

    constructor(address initialOwner, IERC20 _sudaToken)
        Ownable(initialOwner)
    {
        require(address(_sudaToken) != address(0), "escrow: zero token");
        sudaToken = _sudaToken;
    }

    // ────────────────────────────────────────────────────────
    // Lifecycle operations
    // ────────────────────────────────────────────────────────

    /**
     * @notice Create a new quest. Caller must `approve` SUDA spending first.
     *         Reward is pulled into escrow.
     */
    function createQuest(uint256 reward, uint256 deadline)
        external
        whenNotPaused
        nonReentrant
        returns (uint256 questId)
    {
        require(reward > 0, "escrow: zero reward");
        require(deadline > block.timestamp, "escrow: deadline in past");

        require(
            sudaToken.transferFrom(msg.sender, address(this), reward),
            "escrow: SUDA transferFrom failed"
        );

        unchecked {
            questId = ++nextQuestId;
        }
        quests[questId] = Quest({
            creator:  msg.sender,
            assignee: address(0),
            reward:   reward,
            deadline: deadline,
            status:   Status.OPEN
        });

        emit QuestCreated(questId, msg.sender, reward, deadline);
    }

    /**
     * @notice Claim an open quest. Caller becomes the assignee.
     *         First-come-first-served — no second claims.
     */
    function claim(uint256 questId) external whenNotPaused {
        Quest storage q = quests[questId];
        require(q.status == Status.OPEN, "escrow: not open");
        require(block.timestamp < q.deadline, "escrow: past deadline");
        require(msg.sender != q.creator, "escrow: creator cannot claim own quest");

        q.assignee = msg.sender;
        q.status = Status.CLAIMED;

        emit QuestClaimed(questId, msg.sender);
    }

    /**
     * @notice Mark a claimed quest as submitted. Only the assignee can call.
     *         The submission note (description of work) lives off-chain in DB.
     */
    function submit(uint256 questId) external whenNotPaused {
        Quest storage q = quests[questId];
        require(q.status == Status.CLAIMED, "escrow: not claimed");
        require(msg.sender == q.assignee, "escrow: not assignee");

        q.status = Status.SUBMITTED;

        emit QuestSubmitted(questId, msg.sender);
    }

    /**
     * @notice Approve a submitted quest. Only the creator can call.
     *         Sends the reward to the assignee. Terminal state.
     */
    function approve(uint256 questId) external nonReentrant {
        Quest storage q = quests[questId];
        require(q.status == Status.SUBMITTED, "escrow: not submitted");
        require(msg.sender == q.creator, "escrow: not creator");

        address assignee = q.assignee;
        uint256 reward   = q.reward;

        q.status = Status.APPROVED;

        require(
            sudaToken.transfer(assignee, reward),
            "escrow: payout to assignee failed"
        );

        emit QuestApproved(questId, assignee, reward);
    }

    /**
     * @notice Cancel a quest before it's approved. Only the creator can call.
     *         Refunds the reward. Terminal state.
     *         Allowed at any pre-APPROVED state — including SUBMITTED, but in
     *         practice the creator should approve a submitted quest in good faith.
     */
    function cancel(uint256 questId) external nonReentrant {
        Quest storage q = quests[questId];
        require(msg.sender == q.creator, "escrow: not creator");
        require(
            q.status == Status.OPEN ||
            q.status == Status.CLAIMED ||
            q.status == Status.SUBMITTED,
            "escrow: not cancellable"
        );

        uint256 refund = q.reward;
        q.status = Status.CANCELLED;

        require(
            sudaToken.transfer(q.creator, refund),
            "escrow: refund failed"
        );

        emit QuestCancelled(questId, q.creator, refund);
    }

    /**
     * @notice Expire a quest after its deadline. Anyone can call.
     *         Refunds the reward to the creator. Terminal state.
     *         The "anyone can call" pattern means the backend (or any third
     *         party) can clean up overdue quests — escrow correctness doesn't
     *         depend on the creator remembering to act.
     */
    function expire(uint256 questId) external nonReentrant {
        Quest storage q = quests[questId];
        require(
            q.status == Status.OPEN ||
            q.status == Status.CLAIMED ||
            q.status == Status.SUBMITTED,
            "escrow: not expirable"
        );
        require(block.timestamp >= q.deadline, "escrow: not yet expired");

        uint256 refund = q.reward;
        q.status = Status.EXPIRED;

        require(
            sudaToken.transfer(q.creator, refund),
            "escrow: refund failed"
        );

        emit QuestExpired(questId, refund);
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