package tailMail

import (
	"bufio"
	"bytes"
	"log"
	"os"
)

func Tail(fileName string, buffer *bytes.Buffer, oldFileSize int64) (hasNewInfo bool, newFileSize int64, err error) {
	hasNewInfo = false
	newFileSize = oldFileSize

	// 检查文件是否发生变化，
	// 通过比较之前记录的文件尺寸，来判断文件是否被新增内容。
	// 注意，文件被修改，或者删除部分内容，是判断不出来的。
	fileInfo, err := os.Stat(fileName)
	if err != nil {
		log.Println("err:", err)
		return
	}
	newFileSize = fileInfo.Size()

	if newFileSize <= 0 {
		// 文件被清理了。不用发送内容。
		hasNewInfo = false
		err = nil
		return
	}

	tailLen := newFileSize - oldFileSize

	if tailLen != 0 {
		hasNewInfo = true

		// 避免去太长的文本
		if tailLen > 10000 || tailLen < 0 {
			// 太长，截断之， 文件被删除内容，意味着需要重新读取，读取最后的 10000 字符
			tailLen = 10000
		}

		tailLen = tailLen + 100 // 从文件最后往前读的位移量

		tailPos := 0 - tailLen

		err = fileReader(fileName, buffer, tailPos)
	}
	return
}

// 从文件中读取最后 tailLen 长度的内容。
// 第一个不满的空行排除
func fileReader(fileName string, buffer *bytes.Buffer, tailPos int64) (err error) {
	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	log.Println("tailPos：", tailPos)
	file.Seek(tailPos, os.SEEK_END)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {

		buffer.WriteString(scanner.Text())
		buffer.WriteString("\r\n")

		// 获取从文件中读取的内容
		//fmt.Println(scanner.Text())
	}

	return
}
