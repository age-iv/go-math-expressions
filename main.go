package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

func main() {
	// Парсим аргументы командной строки
	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s <input file> <output file>\n", os.Args[0])
		os.Exit(1)
	}
	inputFile := os.Args[1]
	outputFile := os.Args[2]

	// Открываем входной файл для чтения
	in, err := os.Open(inputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка открытия входного файла: %v\n", err)
		os.Exit(1)
	}
	defer in.Close()

	// Создаём или перезаписываем выходной файл
	out, err := os.Create(outputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка создания выходного файла: %v\n", err)
		os.Exit(1)
	}
	defer out.Close()

	// Буферизированная запись
	writer := bufio.NewWriter(out)
	defer writer.Flush()

	// Регулярное выражение для выражений вида "число+число=?" или "число-число=?"
	// Числа могут быть целыми (положительными или отрицательными? по условию, вероятно, только положительные).
	// Для простоты: цифры, затем + или -, затем цифры, затем =?
	re := regexp.MustCompile(`^(\d+)([+\-])(\d+)=\?$`)

	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		line := scanner.Text()
		matches := re.FindStringSubmatch(line)
		if matches == nil {
			continue // строка не подходит, пропускаем
		}
		// matches[1] - первое число, matches[2] - оператор, matches[3] - второе число
		a, _ := strconv.Atoi(matches[1])
		b, _ := strconv.Atoi(matches[3])
		var result int
		switch matches[2] {
		case "+":
			result = a + b
		case "-":
			result = a - b
		}
		// Формируем выходную строку
		outLine := fmt.Sprintf("%s%d\n", line[:len(line)-1], result) // убираем '?' в конце и добавляем результат
		// Записываем в буфер
		if _, err := writer.WriteString(outLine); err != nil {
			fmt.Fprintf(os.Stderr, "Ошибка записи: %v\n", err)
			os.Exit(1)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка чтения входного файла: %v\n", err)
		os.Exit(1)
	}
}
