### 1\. Загальні відомості про лабораторну

**Лабораторна робота №4. Варіант 5. Відлагодження програм на Go.  
Тема: “Розподіл мережевого трафіку між каналами”**

**Мета роботи:**

*   навчитися знаходити та виправляти **compile-time (синтаксичні)** помилки;
*   навчитися запускати програму та порівнювати роботу програми з умовою задачі;
*   навчитися користуватись **debugger’ом у VS Code** (breakpoints, Step Over, Step Into, Watch, Variables, Call Stack);
*   навчитися знаходити та виправляти **логічні (runtime)** помилки;
*   закріпити базовий синтаксис Go: **змінні, слайси, цикли, if, функції, прості структури**.

* * *

### 2\. Постановка задачі (що бачить студент)

Уяви, що є невелика мережа з **трьома каналами зв’язку** (наприклад, три інтернет-канали від різних провайдерів).  
До мережі надходять “пакети” трафіку (умовні значення в Мбіт/с), і система повинна **розподіляти їх між каналами**, щоб:

*   максимально рівномірно розвантажити канали;
*   по можливості **не перевищувати** їхню пропускну здатність (capacity);
*   фіксувати моменти **перевантаження** (overload), коли канал перевищує свою пропускну здатність.

#### Вихідні умови

1.  Є **3 канали** з різною максимальною пропускною здатністю (capacity), наприклад:
    *   `CH-1` — 100 одиниць;
    *   `CH-2` — 120 одиниць;
    *   `CH-3` — 80 одиниць.
2.  Є список “пакетів” трафіку (умовні значення, наприклад: `80, 50, 120, 40, 60, 110`).  
    Кожен елемент цього списку — це трафік, який треба кудись “покласти” (пропустити через один із каналів).
3.  Правильна логіка розподілу:
    *   Перед обробкою кожного пакета система **має вибирати канал з найменшим поточним навантаженням** (`currentLoad`).
    *   Обраний канал отримує весь трафік поточного пакета:  
        `currentLoad = currentLoad + amount`.
    *   Якщо після додавання трафіку `currentLoad` **строго більший за capacity**, вважається, що стався **overload**, і для каналу збільшується лічильник `overloads`.
    *   Також система рахує:
        *   загальну кількість трафіку, що пройшов через всі канали (`totalTraffic`);
        *   загальну кількість подій перевантаження (`overloadEvents`).
4.  Після моделювання програма повинна:
    *   вивести стан кожного каналу:
        *   ім’я каналу;
        *   поточне навантаження (`currentLoad`);
        *   кількість перевантажень (`overloads`);
    *   вивести загальну статистику:
        *   `totalTraffic` — сумарний трафік;
        *   `overloadEvents` — скільки разів хоч якийсь канал перевищував capacity;
        *   канал, який був **найчастіше перевантажений** (за кількістю `overloads`);
        *   текстовий висновок:
            *   якщо перевантажень немає → “Мережа працює стабільно”;
            *   якщо перевантаження були → “Мережа перевантажується, потрібна оптимізація розподілу трафіку”.

⚠️ Опис логіки **правильний**.  
Нижче дано **зламаний код**, який порушує цю логіку й містить **синтаксичні та логічні помилки**.

* * *

### 3\. Зламаний вихідний код програми (`main.go`)

```go
package main

import (
    "fmt"
    "time"
)

type Channel struct {
    name         string
    capacity     int
    currentLoad  int
    overloads    int
}

type NetworkStats struct {
    totalTraffic   int
    overloadEvents int
}

func createChannels() []Channel {
    channels := []Channel{
        {"CH-1", 100, 0, 0},
        {"CH-2", 120, 0, 0},
        {"CH-3", 80, 0, 0},
    }
    return channels
}

func generateTraffic() []int {
    // список "пакетів" трафіку
    traffic := []int{80, 50, 120, 40, 60, 110}
    return traffic
}

// findLeastLoadedChannel повинен знаходити канал з найменшим currentLoad
func findLeastLoadedChannel(channels []Channel) int {
    minIndex := 0
    minLoad := channels[0].currentLoad

    for i := 1; i <= len(channels); i++ {
        if channels[i].currentLoad < minLoad {
            minLoad = channels[i].currentLoad
            minIndex = i
        }
    }

    return minIndex
}

// distributeTraffic розподіляє трафік по каналах
func distributeTraffic(channels []Channel, traffic []int) NetworkStats {
    stats := NetworkStats{}

    fmt.Println("Starting traffic distribution...")
    fmt.Println("Number of channels:", len(channels))
    fmt.Println("Traffic batches:", len(traffic))

    // симулюємо кілька тактів
    for i := 0; i <= len(traffic); i++ {
        amount := traffic[i]
        fmt.Println("\nTick", i+1, "incoming traffic:", amount)

        idx := rand.Intn(len(channels))

        ch := channels[idx]

        // додаємо трафік
        ch.currentLoad = ch.currentLoad + amount

        // рахуємо загальний трафік
        stats.totalTraffic = stats.totalTraffic + amount

        // перевіряємо перевантаження
        if ch.currentLoad >= ch.capacity {
            fmt.Println("WARNING: overload on channel", ch.name)
            stats.overloadEvents = stats.overloadEvents + 1
            ch.overload++
        }

        fmt.Println("Channel", ch.name, "load:", ch.currentLoad, "/", ch.capacity)
    }

    return stats
}

// findMostOverloadedChannel має знаходити канал з найбільшою кількістю overloads
func findMostOverloadedChannel(channels []Channel) Channel {
    most := channels[0]

    for i := 1; i < len(channels); i++ {
        if channels[i].overloads < most.overloads {
            most = channels[i]
        }
    }

    return most
}

func printChannelState(channels []Channel) {
    fmt.Println("\n=== Channels state ===")
    for i := 0; i < len(channels); i++ {
        ch := channels[i]
        fmt.Println(ch.name, "- load:", ch.currentLoad, "/", ch.capacity, "overloads:", ch.overloads)
    }
}

func printSummary(stats NetworkStats, mostOverloaded Channel) {
    fmt.Println("\n=== Network summary ===")
    fmt.Println("Total traffic:", stats.total)
    fmt.Println("Overload events:", stats.overloadEvents)
    fmt.Println("Most overloaded channel:", mostOverloaded.name, "with", mostOverloaded.overloads, "overloads")

    if stats.overloadEvents == 0 {
        fmt.Println("Status: network is stable.")
    } else {
        fmt.Println("Status: network is overloaded, optimization required.")
    }
}

func main() {
    fmt.Println("Network traffic simulation started at", time.Now())

    channels := createChannels()
    traffic := generateTraffic()

    stats := NetworkStats{}

    stats = distributeTraffic(channels, traffic, stats)

    printChannelState(channels)

    mostOverloaded := findMostOverloadedChannel(channels)

    printSummary(stats)
}
```

У коді є:

*   **compile-time помилки** (ти їх побачиш при `go run main.go`).
*   **логічні помилки**.

* * *

### 4\. Інструкції до виконання лабораторної (debug у VS Code)

#### 1\. Підготовка проєкту

1.  Створи папку, наприклад: `lab4_traffic_debug`.
2.  Усередині створити файл `main.go`.
3.  Скопіюй туди весь код із секції 3.
4.  Переконайся, що:
    *   Go встановлено (`go version` у терміналі);
    *   у VS Code встановлено розширення **Go**.

* * *

#### 2\. Виправлення compile-time (синтаксичних) помилок

1.  Відкрий термінал у папці проєкту й запусти:
    ```bash
    go run main.go
    ```
2.  Уважно переглянь перше повідомлення компілятора:
    *   помилки типу `undefined: rand`;
    *   `too many arguments in call to distributeTraffic`;
    *   `stats.total undefined` тощо.
3.  Виправляй помилки **по одній**.
4.  Після кожного виправлення знову запускай:
    ```bash
    go run main.go
    ```
    Переходь далі, коли програма **компілюється** і **запускається**, але працює **не так, як описано в постановці задачі**.

* * *

#### 3\. Налаштування відлагодження у VS Code

1.  Відкрий папку проєкту у VS Code.
2.  Відкрий `main.go`.
3.  Перейди на вкладку **Run and Debug** (ліва панель, ▶️+жук).
4.  Натисни кнопку **“Run and Debug”** → обери конфігурацію **“Go: Launch Package”** (якщо VS Code запитає).
5.  Запам’ятай основні клавіші:
    *   **F5** — запуск відлагодження (**Start Debugging**);
    *   **Shift+F5** — зупинка;
    *   **Ctrl+Shift+F5** — перезапуск.

* * *

#### 4\. Breakpoints (точки зупину)

Постав **breakpoints** у ключових місцях:

1.  У функції `distributeTraffic`:
    *   на рядку з початком циклу `for i := 0; i <= len(traffic); i++ {`;
    *   на рядку, де вибирається канал (там, де `idx := ...`);
    *   на рядку, де обчислюється `ch.currentLoad = ch.currentLoad + amount`;
    *   на рядку з перевіркою перевантаження (`if ch.currentLoad >= ch.capacity { ... }`).
2.  У `findMostOverloadedChannel`:
    *   на рядку з умовою порівняння `if channels[i].overloads < most.overloads {`.
3.  За бажанням — у `main`, перед викликом `distributeTraffic` та перед `printSummary`.

Після цього натисни **F5** і переконайся, що виконання зупиняється на breakpoints.

* * *

#### 5\. Крокування по коду (Step Over / Step Into / Step Out)

Під час зупинки на breakpoint’і використовуй:

*   **F10 (Step Over)** — виконати поточний рядок та перейти до наступного;
*   **F11 (Step Into)** — зайти всередину викликаної функції (наприклад, у `distributeTraffic` або `findMostOverloadedChannel`);
*   **Shift+F11 (Step Out)** — вийти з поточної функції назад у місце виклику.

Рекомендації:

*   простеж, як змінюються:
    *   `amount` (поточний трафік);
    *   `idx` (індекс вибраного каналу);
    *   `ch.currentLoad` і `ch.capacity`;
    *   `stats.totalTraffic` і `stats.overloadEvents`;
*   зверни увагу, чи **дійсно** оновлення `ch.currentLoad` зберігається в `channels[idx]`.

* * *

#### 6\. Використання панелей Variables / Watch / Call Stack

1.  У панелі **Variables** переглядай:
    *   поточні значення `channels`;
    *   змінну `ch`;
    *   структуру `stats`.
2.  У панелі **Watch** додай:
    *   `stats.totalTraffic`;
    *   `stats.overloadEvents`;
    *   `channels[0].currentLoad`;
    *   `channels[1].currentLoad`;
    *   `channels[2].currentLoad`;
    *   `channels[0].overloads`, `channels[1].overloads`, `channels[2].overloads`.
3.  У **Call Stack** подивись:
    *   як виконання заходить у `distributeTraffic` з `main`;
    *   як пізніше викликається `findMostOverloadedChannel`.

Це допоможе бачити, **де саме** в коді ти знаходишся й хто кого викликав.

* * *

#### 7\. Пошук і виправлення логічних помилок

Треба знайти й виправити **мінімум 5 логічних помилок**.

Для кожної знайденої логічної помилки в звіті:

*   наведи фрагмент **“Було”** (до виправлення);
*   фрагмент **“Стало”** (після виправлення);
*   дай коротке пояснення:  
    _“Було: цикл доходив до `len(traffic)`, тому виникав вихід за межі слайса. Стало: замінив `<=` на `<`, тепер індекси коректні.”_

* * *

#### 8\. Перевірка результату

Після виправлення **усіх синтаксичних помилок** і основних **логічних помилок**:

1.  Запусти програму:
    ```bash
    go run main.go
    ```
2.  Перевір, що:
    *   всі “пакети” трафіку з `generateTraffic` обробляються;
    *   канали отримують трафік згідно з логікою **найменш завантаженого** каналу;
    *   події перевантаження рахуються коректно (тільки коли `currentLoad` перевищує `capacity`);
    *   `totalTraffic` дорівнює сумі всіх значень у масиві трафіку;
    *   функція `findMostOverloadedChannel` справді повертає канал з найбільшою кількістю `overloads`;
    *   текстовий висновок (“stable” / “overloaded”) відповідає кількості подій `overloadEvents`.
3.  Зафіксуй у звіті:
    *   список знайдених та виправлених **compile-time** помилок;
    *   список знайдених та виправлених **логічних** помилок;
    *   фінальний скрін/копію виводу програми.

* * *