

result = ""
if (bc.TargetGroup.DeckInBoard.length < ca.Position + 1) { // hit hero
    bc.TargetGroup.Hero.CHP -= 150
    result = "User" + bc.CurrentGroup.UserID + " deal 150 dg to " + bc.TargetGroup.Hero.Name
} else {
    bc.TargetGroup.DeckInBoard[ca.Position].CHP -= 150
    result = "User" + bc.CurrentGroup.UserID + " deal 150 dg to " + bc.TargetGroup.DeckInBoard[ca.Position].Name
}
result