# REST API BACKEND

## GetUser
Получить общий скор и информацию о еще не переработанном мусоре

`GET /bihack/rest/users/{firebase_id_or_whatever}`

*Создает нового пользователя если такой айди не был зарегестрирован ранее*

    {
      "score": 10,
      "in_progress": {
          "inbin": 3,
          "processing": 2,
      }
    }

## GetItem
Получить тип товара по бар коду и добавить его в бд если это первый окуранс

`POST /bihack/rest/item/[bar-code]/`
    
    {
        "type": "plastic"
    }


## AddHistoryRecord
Добавить новую запись о сдаче мусора

`POST /bihack/rest/history/`

### Request:

    {
        "user": "ABC",
        "residence": "ExpoPlaza",
        "type": "plastic",
        "amount": 10
    }

### Response:
200 on success 500 on error

## AddResidence
Добавить новый жк в дб

`POST /bihack/rest/new_residential/`

### Request:

    {
        "name": "ExpoBoulevard",
        "coordinates": "1.1.1.1"
    }

### Response:
200 on success 500 on error

## GetGPSCoordinates
Отправить координаты жучка

`GET /bihack/rest/coords/{coordinate}`

200 on success 500 on error