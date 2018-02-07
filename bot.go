//ввод кораблей компьютера
import "math/rand"

func BotInput() {
	for k := 0; k < lengthOfShip; k++ { //цикл по длине корабля
		for l := lengthOfShip; l > k; l-- { //цикл по кол-ву кораблей
			flag := true
			var n1, n2, l1, l2 uint8
			for {
				flag = false
				l1 = uint8(rand.Int()%10) + 1
				n1 = uint8(rand.Int()%10) + 1
				l2 = uint8(rand.Int() % 2)
				if l2 == 1 {
					l2 = l1 + uint8(k)
					n2 = n1
				} else {
					n2 = n1 + uint8(k)
					l2 = l1
				}
				//проверяем что наш корабль не пересекается с уже заданными
				for f := l1 - 1; f <= l2+1; f++ {
					for g := n1 - 1; g <= n2+1; g++ {
						if computer[g][f] == 1 {
							flag = true
						}
					}
				}
				if ((flag == true) || (n2 > 10) || (l2 > 10)) == false {
					break
				}
			}
			//если мы тут, значит корабль успешно введён
			//заносим его в наш массив(поле) и помечаем соседние с кораблём клетки
			for i := l1; i <= l2; i++ {
				for j := n1; j <= n2; j++ {
					computer[j][i] = 1
				}
			}
			for f := l1 - 1; f <= l2+1; f++ {
				for g := n1 - 1; g <= n2+1; g++ {
					if computer[g][f] != 1 {
						computer[g][f] = 8
					}
				}
			}
		}
	}
}