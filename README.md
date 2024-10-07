Comming soon!


for migration we should create folder (usualy is migrations)
then we should create file with the follog format id_action_table_name.up.sql (1481574547_create_users_table.up.sql)


Migration command
docker run -v /Users/jovid/Desktop/Dev/microCmp/migrations:/migrations --network host migrate/migrate -path=/migrations/ -database "postgresql://postgres:test@localhost:5432/postgres?sslmode=disable" up

--Windows
docker run -v C:/Users/dzhovid.nurov/Desktop/dev/MicroCmp/migrations:/migrations --network host migrate/migrate -path=/migrations/ -database "postgresql://postgres:test@localhost:5432/postgres?sslmode=disable" up

docker-compose exec web go run scripts/seed.go




"availableServs": [
        {
            "groupID": 0,
            "unitID": 3,
            "id": 4116,
            "name": "Мегабайты за звонки",
            "desc": "Получайте по 10 МБ за каждую минуту входящего звонка!",
            "price": 2,
            "payPeriod": "/день"
        },
        {
            "groupID": 0,
            "unitID": 3,
            "id": 3810,
            "name": "Фиолетовый день",
            "desc": "1000 МБ трафика сроком на 1 сутки",
            "price": 5,
            "payPeriod": ""
        },
        {
            "groupID": 0,
            "unitID": 4,
            "id": 2245,
            "name": "Simфония",
            "desc": "Замени обыкновенный гудок на любимые мелодии",
            "price": 0.25,
            "payPeriod": "/день"
        },
        {
            "groupID": 0,
            "unitID": 4,
            "id": 4757,
            "name": "Пресса",
            "desc": "Здесь собраны лицензионные версии лучших журналов, популярных во всем мире!",
            "price": 0.75,
            "payPeriod": "/день"
        },
        {
            "groupID": 0,
            "unitID": 4,
            "id": 4758,
            "name": "Караоке",
            "desc": "Cервис позволяет петь любимые песни с использованием только мобильного телефона!",
            "price": 1.43,
            "payPeriod": "/день"
        },
        {
            "groupID": 0,
            "unitID": 4,
            "id": 4760,
            "name": "USSD Викторина 'Mоментальная'",
            "desc": "Отвечай на вопросы по USSD и выиграй до 2000 сомони в месяц!",
            "price": 0.6,
            "payPeriod": "/день"
        },
        {
            "groupID": 0,
            "unitID": 4,
            "id": 4761,
            "name": "SMS Викторина 'Единая'",
            "desc": "Отвечайте на вопросы и выиграйте денежные призы!",
            "price": 0.5,
            "payPeriod": "/день"
        },
        {
            "groupID": 0,
            "unitID": 4,
            "id": 4762,
            "name": "Видео-портал",
            "desc": "Смотрите фильмы, мультфильмы, сериалы и ТВ каналы в Full HD и без рекламы на портале video.tcell.tj!",
            "price": 1.3,
            "payPeriod": "/день"
        },
        {
            "groupID": 0,
            "unitID": 4,
            "id": 4765,
            "name": "Знакомства 6900",
            "desc": "Заходите на pinme.tcell.tj знакомьтесь, общайтесь, дружите!",
            "price": 1,
            "payPeriod": "/день"
        },
        {
            "groupID": 0,
            "unitID": 4,
            "id": 4766,
            "name": "Игровой портал",
            "desc": "Играйте в своем телефоне в разные игры на портале games.tcell.tj",
            "price": 1,
            "payPeriod": "/день"
        },
        {
            "groupID": 0,
            "unitID": 4,
            "id": 4773,
            "name": "Секреты домохозяек",
            "desc": "Узнайте разные секреты и доставайтесь надежной домохозяйкой!",
            "price": 0.7,
            "payPeriod": "/день"
        },
        {
            "groupID": 0,
            "unitID": 4,
            "id": 4775,
            "name": "Ошикона",
            "desc": "Звоните бесплатно на 5050 знакомьтесь, общайтесь, дружите!",
            "price": 0.8,
            "payPeriod": "/день"
        },
        {
            "groupID": 0,
            "unitID": 4,
            "id": 4776,
            "name": "Хандаборон",
            "desc": "Звоните на бесплатный номер 9402 и разыграйте своих друзей!",
            "price": 0.8,
            "payPeriod": "/день"
        },
        {
            "groupID": 0,
            "unitID": 4,
            "id": 4802,
            "name": "Популярные мультсериалы",
            "desc": "Популярные мультики в вашем телефоне и для вашего ребенка!",
            "price": 0.8,
            "payPeriod": "/день"
        },
        {
            "groupID": 0,
            "unitID": 4,
            "id": 4803,
            "name": "Журналы",
            "desc": "Портал свежих журналов! Подключайся и листай!",
            "price": 0.6,
            "payPeriod": "/день"
        },
        {
            "groupID": 0,
            "unitID": 4,
            "id": 4804,
            "name": "Вселенная Angry Birds",
            "desc": "Мультики от Angry birds! Смотри и насладись!",
            "price": 0.8,
            "payPeriod": "/день"
        },
        {
            "groupID": 0,
            "unitID": 4,
            "id": 4805,
            "name": "Сериалы и ТВ",
            "desc": "Услуга для любителей сериалов и ТВ шоу",
            "price": 1,
            "payPeriod": "/день"
        },
        {
            "groupID": 0,
            "unitID": 4,
            "id": 4807,
            "name": "Автопортал",
            "desc": "Автопортал, новости про авто в твоем телефоне!",
            "price": 1,
            "payPeriod": "/день"
        },
        {
            "groupID": 0,
            "unitID": 4,
            "id": 4808,
            "name": "Киберспортивный",
            "desc": "Все о мире киберспорта, подпишись и узнай новости про мир игр!",
            "price": 1,
            "payPeriod": "/день"
        },
        {
            "groupID": 0,
            "unitID": 4,
            "id": 4810,
            "name": "Подкасты",
            "desc": "Подкасты на web! Слушайте интересные рассказы!",
            "price": 1,
            "payPeriod": "/день"
        },
        {
            "groupID": 0,
            "unitID": 4,
            "id": 4813,
            "name": "Фанбокс Плеер",
            "desc": "Слушайте шутки и приколы!",
            "price": 0.5,
            "payPeriod": "/день"
        },
        {
            "groupID": 0,
            "unitID": 4,
            "id": 4814,
            "name": "Фанбокс Анекдоты",
            "desc": "Смешные анекдоты в твоем телефоне!",
            "price": 0.5,
            "payPeriod": "/день"
        },
        {
            "groupID": 0,
            "unitID": 4,
            "id": 4815,
            "name": "Фанбокс Знакомства",
            "desc": "Новые друзья новые знакомства! Подпишись и вперед!",
            "price": 0.5,
            "payPeriod": "/день"
        },
        {
            "groupID": 0,
            "unitID": 4,
            "id": 4816,
            "name": "Фанбокс Тесты",
            "desc": "Тесты в твоем телефоне!",
            "price": 0.5,
            "payPeriod": "/день"
        },
        {
            "groupID": 0,
            "unitID": 4,
            "id": 114,
            "name": "АОН",
            "desc": "Определи скрытые номера!",
            "price": 0.22,
            "payPeriod": "/день"
        },
        {
            "groupID": 0,
            "unitID": 4,
            "id": 4774,
            "name": "Хушнаво 0888",
            "desc": "Интересные рассказы и истории на портале 0888, слушайте и наслаждайтесь!",
            "price": 0.8,
            "payPeriod": "/день"
        },
        {
            "groupID": 0,
            "unitID": 4,
            "id": 1782,
            "name": "Викторина меломания",
            "desc": "Выиграй один из 3 смартфонов самсунг за месяц!",
            "price": 1,
            "payPeriod": "/день"
        },
        {
            "groupID": 0,
            "unitID": 4,
            "id": 2807,
            "name": "Oila TV (app)",
            "desc": "Любимые фильмы и ТВ передачи в твоем телефоне!",
            "price": 1.8,
            "payPeriod": "/день"
        },
        {
            "groupID": 0,
            "unitID": 4,
            "id": 4763,
            "name": "Викторина 'Мегабайты за знания'",
            "desc": "Участвуйте и выиграйте до 5000 мегабайтов!",
            "price": 1,
            "payPeriod": "/день"
        },
        {
            "groupID": 0,
            "unitID": 1,
            "id": 3936,
            "name": "+50 Минут",
            "desc": "На другие операторы по РТ",
            "price": 10,
            "payPeriod": ""
        },
        {
            "groupID": 0,
            "unitID": 1,
            "id": 3939,
            "name": "+300 Минут",
            "desc": "На другие операторы по РТ",
            "price": 55,
            "payPeriod": ""
        },
        {
            "groupID": 0,
            "unitID": 4,
            "id": 4777,
            "name": "Этикет в гостях",
            "desc": "Узнай как действовать, вести себя в гостях, на улице или в обществе!",
            "price": 0.7,
            "payPeriod": "/день"
        },
        {
            "groupID": 0,
            "unitID": 4,
            "id": 4806,
            "name": "Уроки",
            "desc": "Видео уроки в твоем телефоне! Научиться играть на гитаре, готовить и много другое!",
            "price": 0.7,
            "payPeriod": "/день"
        },
        {
            "groupID": 0,
            "unitID": 4,
            "id": 4809,
            "name": "Спортивный портал",
            "desc": "Новости о мире спорта! Новые результаты! Будь в тренде мировых новостей о спорте, узнай все в своем телефоне!",
            "price": 1,
            "payPeriod": "/день"
        },
        {
            "groupID": 0,
            "unitID": 4,
            "id": 4772,
            "name": "Мудрые мысли",
            "desc": "Ежедневные SMS с мудрыми мыслями и советами от философов и мудрецов всех времен!",
            "price": 0.7,
            "payPeriod": "/день"
        },
        {
            "groupID": 0,
            "unitID": 4,
            "id": 4136,
            "name": "Смартбук",
            "desc": "Книги на любой вкус! Читайте любимые книги в своем телефоне всего за 1 сомони в день.",
            "price": 1,
            "payPeriod": "/день"
        },
        {
            "groupID": 0,
            "unitID": 1,
            "id": 3937,
            "name": "+100 Минут",
            "desc": "На другие операторы по РТ",
            "price": 20,
            "payPeriod": ""
        },
        {
            "groupID": 0,
            "unitID": 1,
            "id": 3938,
            "name": "+200 Минут",
            "desc": "На другие операторы по РТ",
            "price": 35,
            "payPeriod": ""
        },
        {
            "groupID": 0,
            "unitID": 4,
            "id": 4684,
            "name": "Ночной безлимит для линейки Салом",
            "desc": "",
            "price": 10,
            "payPeriod": ""
        },
        {
            "groupID": 0,
            "unitID": 4,
            "id": 4759,
            "name": "Игра Мореплаватель",
            "desc": "Играй и выиграй деньги на свой баланс!",
            "price": 1,
            "payPeriod": "/день"
        },
        {
            "groupID": 0,
            "unitID": 4,
            "id": 4075,
            "name": "WEB-игра 2048",
            "desc": "Играйте в популярное игре 2048 и выиграйте денежные призы!",
            "price": 0.9,
            "payPeriod": "/день"
        }
    ]