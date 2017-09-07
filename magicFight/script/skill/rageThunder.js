// rageThunder skill 

function inArray(target, arr) {
    for (var i=0;i<arr.length;i++) {
        if (arr[i] == target) return true
    }
    return false
}

result = [];
if (bc.TargetGroup.DeckInBoard.length <= 3) {
    for (var i in bc.TargetGroup.DeckInBoard) {
        bc.TargetGroup.DeckInBoard[i].CHP -= 150
        result.push("User" + bc.CurrentGroup.UserID + " deal 150 dg to " + bc.TargetGroup.DeckInBoard[i].Name)
    }
} else {
    var countList = [];
    for (var i in bc.TargetGroup.DeckInBoard) {
        if (countList.length >= 3) break;
        if (!inArray(i, countList)) {
            if (Math.random > 0.5) {
                countList.push(i)
                bc.TargetGroup.DeckInBoard[i].CHP -= 150
                result.push("User" + bc.CurrentGroup.UserID + " deal 150 dg to " + bc.TargetGroup.DeckInBoard[i].Name)
            }
        }
    }
}
result = result.join("\n")
result
