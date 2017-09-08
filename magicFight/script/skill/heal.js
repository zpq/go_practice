function inArray(target, arr) {
    for (var i=0;i<arr.length;i++) {
        if (arr[i] == target) return true
    }
    return false
}

result = [];

ca.CHP += 300
if (ca.CHP > ca.HP) {
    ca.CHP = ca.HP
}
result.push("User" + bc.CurrentGroup.UserID + " add 300 HP to self")

var l = bc.CurrentGroup.DeckInBoard.length;

if (l <= 3) {
    for (var i in bc.CurrentGroup.DeckInBoard) {
        bc.CurrentGroup.DeckInBoard[i].CHP += 200
        result.push("User" + bc.CurrentGroup.UserID + " add 200 HP to " + bc.CurrentGroup.DeckInBoard[i].Name)
    }
} else {
    var countList = [];
    for (var i in bc.CurrentGroup.DeckInBoard) {
        if (countList.length >= 3) break;
        if (!inArray(i, countList)) {
            if (Math.random > 0.5) {
                countList.push(i)
                bc.CurrentGroup.DeckInBoard[i].CHP += 200
                result.push("User" + bc.CurrentGroup.UserID + " add 200 HP to " + bc.CurrentGroup.DeckInBoard[i].Name)
            }
        }
    }
}

result = result.join("\n")
result
