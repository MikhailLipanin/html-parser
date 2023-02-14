# html-parser

Приложение, которое парсит таблицу с данными в указанном HTML и заполняет эти данные в Google Таблицу.

В приложении используются пакеты:

- [viper](github.com/spf13/viper) для работы с конфигурационными файлами
- [spreadsheet.v2](gopkg.in/Iwark/spreadsheet.v2) для взаимодействия с Google Doc API
- [colly](github.com/gocolly/colly/v2) для парсинга HTML файлов