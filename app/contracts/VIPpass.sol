// SPDX-License-Identifier: MIT
pragma solidity ^0.8.28;

import "@openzeppelin/contracts/token/ERC721/ERC721.sol";
import "@openzeppelin/contracts/token/ERC721/extensions/ERC721Enumerable.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

contract VIPPass is ERC721, ERC721Enumerable, Ownable {
    uint256 private _nextId;

    constructor() ERC721("VIPPass", "VIP") Ownable(msg.sender) {}

    function mint(address to) external onlyOwner returns (uint256) {
        require(balanceOf(to) == 0, "Already holds a pass");
        uint256 id = ++_nextId;
        _safeMint(to, id);
        return id;
    }

    function burn(address from) external onlyOwner {
        require(balanceOf(from) == 1, "No pass to burn");
        uint256 tokenId = tokenOfOwnerByIndex(from, 0);
        _burn(tokenId);
    }

    function _update(
        address to,
        uint256 tokenId,
        address auth
    ) internal override(ERC721, ERC721Enumerable) returns (address) {
        address from = _ownerOf(tokenId);
        require(from == address(0) || to == address(0), "Non-transferable");
        return super._update(to, tokenId, auth);
    }

    function _increaseBalance(
        address account,
        uint128 value
    ) internal override(ERC721, ERC721Enumerable) {
        super._increaseBalance(account, value);
    }

    function supportsInterface(
        bytes4 interfaceId
    ) public view override(ERC721, ERC721Enumerable) returns (bool) {
        return super.supportsInterface(interfaceId);
    }
}
