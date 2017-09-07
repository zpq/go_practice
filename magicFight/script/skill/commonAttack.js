
result = ""
console.log("position " + ca.Position)
if (bc.TargetGroup.DeckInBoard.length < ca.Position + 1) { // 对应位置没有卡牌，则攻击英雄
    bc.TargetGroup.Hero.CHP -= ca.CAttack
    result = "User" + bc.CurrentGroup.UserID + " deal  "+ca.CAttack+ "dg to " + bc.TargetGroup.Hero.Name
} else {
    bc.TargetGroup.DeckInBoard[ca.Position].CHP -= ca.CAttack
    result = "User" + bc.CurrentGroup.UserID + " deal "+ca.CAttack+" dg to " + bc.TargetGroup.DeckInBoard[ca.Position].Name
}
result