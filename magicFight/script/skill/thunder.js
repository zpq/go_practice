// thunder skill

result = "";
if (ca.ActorType == 1) {
    bc.TargetGroup.Hero.CHP -= 300;
    result = "User" + bc.CurrentGroup.UserID + " deal 300 dg to " + bc.TargetGroup.Hero.Name
} else {
    if (bc.TargetGroup.DeckInBoard[ca.Position]) {
        bc.TargetGroup.DeckInBoard[ca.Position].CHP -= 300
        result = "User" + bc.CurrentGroup.UserID + " deal 300 dg to " + bc.TargetGroup.DeckInBoard[ca.Position].Name
    } else {
        bc.TargetGroup.Hero.CHP -= 300;
        result = "User" + bc.CurrentGroup.UserID + " deal 300 dg to " + bc.TargetGroup.Hero.Name
    }
}
result