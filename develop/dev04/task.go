/*
=== Поиск анаграмм по словарю ===

Напишите функцию поиска всех множеств анаграмм по словарю.
Например:
'пятак', 'пятка' и 'тяпка' - принадлежат одному множеству,
'листок', 'слиток' и 'столик' - другому.

Входные данные для функции: ссылка на массив - каждый элемент которого - слово на русском языке в кодировке utf8.
Выходные данные: Ссылка на мапу множеств анаграмм.
Ключ - первое встретившееся в словаре слово из множества
Значение - ссылка на массив, каждый элемент которого, слово из множества. Массив должен быть отсортирован по возрастанию.
Множества из одного элемента не должны попасть в результат.
Все слова должны быть приведены к нижнему регистру.
В результате каждое слово должно встречаться только один раз.

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/
package main

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
)

func anagram(data []string) map[string][]string {
	res := make(map[string][]string)
	exists := make(map[string]struct{})
	for i, v := range data {
		exists[v] = struct{}{}

		for j := i + 1; j < len(data); j++ {
			w := strings.ToLower(strings.TrimSpace(data[j]))
			if _, ok := exists[w]; ok {
				continue
			}
			if areAnagrams(v, w) {
				res[v] = append(res[v], w)
				exists[w] = struct{}{}
			}
		}
	}

	for i, v := range res {
		sort.Strings(v)
		res[i] = v
	}

	return res
}

func frequencyMap(s string) map[rune]int {
	freqMap := make(map[rune]int)
	for _, r := range s {
		freqMap[r]++
	}
	return freqMap
}

func areAnagrams(s1, s2 string) bool {
	if len(s1) != len(s2) {
		return false
	}
	return reflect.DeepEqual(frequencyMap(s1), frequencyMap(s2))
}

func main() {
	data := []string{"пятак", "пятка", "тяпка", "листок", "слиток", "СтоЛИК", "ПЯТКА", "листок"}
	res := anagram(data)

	fmt.Println(res)
}
