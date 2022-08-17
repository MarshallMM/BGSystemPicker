# BGSystemPicker
Board game system picker discord bot to get the group picking games

Commands:
!pick <game>
  -adds <game> to pick list
!veto <game> 
  -adds <game> to veto list
!list
  -lists out the current picks and vetos
!roll
  -sudo ramdomly pick a game out of the pick list that is not vetoed. 
  -A string of the current date and games is put through a hash function, and a number is generated to be the picked game.
  -!roll will return the same game every time if the pick list, veto list, and date are the same. 
  -This is to prevent randomly rolling until a desired choice is made.
  
