package telegramEvents

const StartMessage = `
Привет! Я бот менеджер паролей.
Я был создан в качестве тестового задания для компании VK.
Если тебе интересно, как я работаю, то можешь посмотреть мой код на GitHub 🌐: github.com/Phaseant/PasswordManagerBot
Если ты хочешь узнать, что я умею, то введи /help
Сделано на Go с ❤️`

const HelpMessage = `
/start - посмотреть приветственное сообщение
/help - получить справку о командах
/add - добавить запись
/delete - удалить запись
/get - посмотреть все записи сервиса

Также ты можешь пользоваться меню.`

const unknownMessage = `Я не знаю такой команды 😥
Введи /help, чтобы узнать, что я умею.`

const AddErrorMessage = `Произошла ошибка при добавлении записи 😥
Попробуйте еще раз`

const DeleteErrorMessage = `Произошла ошибка при удалении записи 😥
Попробуйте еще раз`

const GetErrorMessage = `Произошла ошибка при получении записи 😥
Попробуйте еще раз`

const AddSuccessMessage = `Запись успешно добавлена! 🎉`

const DeleteSuccessMessage = `Запись успешно удалена! 🎉`

const NotFoundMessage = `Запись не найдена 😥`
