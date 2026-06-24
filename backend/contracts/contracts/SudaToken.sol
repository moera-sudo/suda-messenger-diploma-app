pragma solidity 0.8.28;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "@openzeppelin/contracts/token/ERC20/extensions/ERC20Pausable.sol";
import "@openzeppelin/contracts/access/Ownable.sol";


contract SudaToken is ERC20, ERC20Pausable, Ownable {
    uint256 public constant INITIAL_SUPPLY = 100_000_000 * 10 ** 18;


    constructor(address initialOwner)
        ERC20("Suda", "SUDA")
        Ownable(initialOwner)
    {
        _mint(initialOwner, INITIAL_SUPPLY);
    }

    /// @notice Mint additional SUDA. Only the owner (treasury) can call this.
    function mint(address to, uint256 amount) external onlyOwner {
        _mint(to, amount);
    }

    /// @notice Pause all token transfers. Only owner.
    function pause() external onlyOwner {
        _pause();
    }

    /// @notice Resume token transfers. Only owner.
    function unpause() external onlyOwner {
        _unpause();
    }

    function _update(address from, address to, uint256 value)
        internal
        override(ERC20, ERC20Pausable)
    {
        super._update(from, to, value);
    }
}