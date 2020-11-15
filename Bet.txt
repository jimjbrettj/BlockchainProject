pragma solidity >=0.4.22 <0.7.0;

/**
 * @title Bet
 * @dev Participate in a betting application
 */
contract Bet {
    address payable public owner;
    uint256 public minimumBet;
    uint256 public totalBet;
    uint256 public cap;
    address payable[] public players;
    
    struct Player {
       uint256 amountBet;
    }
	
	// maps player addres to playerBet   
    mapping(address => Player) public playerBet;
    receive() external payable {}
    
    //constructs contract
    constructor() public payable {
        owner = msg.sender;
        minimumBet = 1;
		cap = 5;
		totalBet = 0;
	}
	
	//allows owner to kill contract
	function kill() public {
      if(msg.sender == owner) selfdestruct(owner);
    } 

    //see if player has already bet
    function checkPlayerExists(address player) public view returns(bool){
      for(uint256 i = 0; i < players.length; i++){
         if(players[i] == player) return true;
      }
      return false;
    }
   
    function random() public view returns(uint) {
      //return uint(keccak256(block.timestamp, block.difficulty))%5;
      return uint(keccak256(abi.encode(block.difficulty, block.timestamp)))%5;
    }
    
    //reset contract state
    function reset() private {
        totalBet = 0;
        for(uint256 i = 0; i < players.length; i++){
          delete playerBet[players[i]];
        }
        delete players;
    }
    
    //built for testing
    function getNumPlayers() public view returns(uint){
        return players.length;
    }
   
    //place bet
    function bet() public payable returns(bool) {
      require(!checkPlayerExists(msg.sender));
      require(msg.value >= minimumBet);
      
      playerBet[msg.sender].amountBet = msg.value;
      players.push(msg.sender);
      totalBet += msg.value;
      
      //if 5 people have bet then pick winner and reset
      if(players.length >= cap) {
          //pick winner
          uint w = random(); 
          //transfer winnings to winner
          address payable winner = players[w];
          if(!winner.send(totalBet)) {
              return false;
          }
          //reset the bet
          reset();
      }
      return true;
    }
}