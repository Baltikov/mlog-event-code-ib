package mlog

import (
	"errors"
	"fmt"
)

// Standard levels
var (
	LvlPanic = LogLevel{ID: 0, Name: "panic", Stacktrace: true}
	LvlFatal = LogLevel{ID: 1, Name: "fatal", Stacktrace: true}
	LvlError = LogLevel{ID: 2, Name: "error"}
	LvlWarn  = LogLevel{ID: 3, Name: "warn"}
	LvlInfo  = LogLevel{ID: 4, Name: "info"}
	LvlDebug = LogLevel{ID: 5, Name: "debug"}
	LvlTrace = LogLevel{ID: 6, Name: "trace"}
	// used by redirected standard logger
	LvlStdLog = LogLevel{ID: 10, Name: "stdlog"}
	// used only by the logger
	LvlLogError = LogLevel{ID: 11, Name: "logerror", Stacktrace: true}
)

// Register custom (discrete) levels here.
// !!!!! ID's must not exceed 32,768 !!!!!!
var (
	// used by the audit system
	LvlAuditAPI     = LogLevel{ID: 100, Name: "audit-api"}
	LvlAuditContent = LogLevel{ID: 101, Name: "audit-content"}
	LvlAuditPerms   = LogLevel{ID: 102, Name: "audit-permissions"}
	LvlAuditCLI     = LogLevel{ID: 103, Name: "audit-cli"}

	// used by the TCP log target
	LvlTcpLogTarget = LogLevel{ID: 120, Name: "TcpLogTarget"}

	// used by Remote Cluster Service
	LvlRemoteClusterServiceDebug = LogLevel{ID: 130, Name: "RemoteClusterServiceDebug"}
	LvlRemoteClusterServiceError = LogLevel{ID: 131, Name: "RemoteClusterServiceError"}
	LvlRemoteClusterServiceWarn  = LogLevel{ID: 132, Name: "RemoteClusterServiceWarn"}

	// used by Shared Channel Sync Service
	LvlSharedChannelServiceDebug            = LogLevel{ID: 200, Name: "SharedChannelServiceDebug"}
	LvlSharedChannelServiceError            = LogLevel{ID: 201, Name: "SharedChannelServiceError"}
	LvlSharedChannelServiceWarn             = LogLevel{ID: 202, Name: "SharedChannelServiceWarn"}
	LvlSharedChannelServiceMessagesInbound  = LogLevel{ID: 203, Name: "SharedChannelServiceMsgInbound"}
	LvlSharedChannelServiceMessagesOutbound = LogLevel{ID: 204, Name: "SharedChannelServiceMsgOutbound"}

	// add more here ...
)

type EventCodeGroup struct {
	AppStarted         int
	UserLoggedIn       int
	DbQuery            int
	TaskCompleted      int
	ConfigChanged      int
	ServiceStarted     int
	ServiceStopped     int
	OperationCompleted int
	LdapConnectionSuccess int
    LdapAuthSuccess int
	UserIdentified int
	LdapGroupIdentified int
	UserRoleAssigned int
	ServerStarting int
	RateLimitingEnabled int
	ServerListening int
	ServerStopped int
	
}

var INFO = EventCodeGroup{
	AppStarted:         1001, // Успешный запуск программы
	UserLoggedIn:       1002, // Вход пользователя в систему
	DbQuery:            1003, // Запрос к базе данных
	TaskCompleted:      1004, // Успешное выполнение задачи
	ConfigChanged:      1005, // Изменения конфигурации
	ServiceStarted:     1006, // Запуск сервиса
	ServiceStopped:     1007, // Остановка сервиса
	OperationCompleted: 1008, // Завершение операции
	LdapConnectionSuccess: 1009, // Успешное подключение к LDAP серверу ib
	LdapAuthSuccess: 1010, // Успешная аутентификация через LDAP ib
	UserIdentified: 1011, // Пользователь успешно идентифицирован по заголовку
	LdapGroupIdentified: 1012, // Успешно идентифицирована группа пользователя через LDAP ib
	UserRoleAssigned: 1014, // Назначение роли пользователю.
	ServerStarting: 1015, //  начало процесса запуска сервера.
	RateLimitingEnabled: 1016, // ограничение скорости запросов включено.
	ServerListening: 1017, //  успешное начало прослушивания на указанном адресе
	ServerStopped: 1018, // успешная остановку сервера.
}

type WarningCodeGroup struct {
	SuspiciousActivity int
	SlowQuery          int
	LowResource        int
	ParsingError       int
	UnsupportedConfig  int
	HighApiLatency     int
	NotFound           int
	DbRecycleStart     int
	DbRecycleEnd       int
	ServerShutdownWarning int

}

var WARN = WarningCodeGroup{
	SuspiciousActivity: 2001, // Подозрительное поведение пользователя
	SlowQuery:          2002, // Длительное выполнение запроса
	LowResource:        2003, // Низкий уровень ресурсов (память, диск)
	ParsingError:       2004, // Ошибки парсинга
	UnsupportedConfig:  2005, // Неподдерживаемая конфигурация
	HighApiLatency:     2006, // Повышенное время отклика внешнего API
	NotFound:           2007, // Обработка ошибки 404 - маршрут не найден
	DbRecycleStart: 2008, // Начало процесса перезапуска подключения к базе данных
	DbRecycleEnd: 2009, // Завершение процесса перезапуска подключения к базе данных
	ServerShutdownWarning: 2010, // предупреждение, если возникает ошибка при попытке корректно остановить сервер
}

type ErrorCodeGroup struct {
	DbConnectionFailed    int
	QueryExecutionFailed  int
	IOError               int
	AuthFailed            int
	DataIntegrationFailed int
	LdapAuthFailed int
	LdapGroupSearchFailed int
	LdapGroupNotFound int
	TemplateParseFailed   int
	TemplateExecutionFailed int
	ConfigWatcherInitFailed int
	RateLimitExceeded int




}

var ERROR = ErrorCodeGroup{
	DbConnectionFailed:    3001, // Ошибка подключения к базе данных
	QueryExecutionFailed:  3002, // Ошибка выполнения запроса
	IOError:               3003, // Ошибка ввода/вывода
	AuthFailed:            3004, // Сбой при аутентификации пользователя
	DataIntegrationFailed: 3005, // Неуспешная интеграция данных
	LdapAuthFailed: 3006, // Ошибка аутентификации через LDAP ib
	LdapGroupSearchFailed: 3007, // Ошибка при поиске групп пользователя в LDAP ib
	LdapGroupNotFound: 3008, // Группы пользователя не найдены в LDAP ib
	TemplateParseFailed: 3009, // Ошибка при разборе HTML-шаблона
	TemplateExecutionFailed: 3010, // Ошибка при выполнении HTML-шаблона
	ConfigWatcherInitFailed: 3011, // Ошибка инициализации наблюдателя конфигурации
	RateLimitExceeded: 3012, // Превышено ограничение скорости запросов
}

type CriticalFatalCodeGroup struct {
	CriticalServiceFailed int
	DataLoss              int
	ImmediateAttention    int
	AppCrash              int
	BusinessProcessFailed int
	RateLimitingFailed    int
	ServerListenFailed    int
	ServerServeFailed     int
	IndexOperationFailed  int

}

var CRITICAL = CriticalFatalCodeGroup{
	CriticalServiceFailed: 4001, // Сбой критической службы
	DataLoss:              4002, // Потеря данных
	ImmediateAttention:    4003, // Ошибка, требующая немедленного внимания
	AppCrash:              4004, // Аварийное завершение приложения
	BusinessProcessFailed: 4005, // Поломка важного бизнес-процесса
	RateLimitingFailed: 4006, // Ошибка при применении ограничения скорости
	ServerListenFailed: 4007, // Ошибка при попытке начать прослушивание на заданном адрес
	ServerServeFailed: 4008, // Ошибка при попытке запустить сервер.
	IndexOperationFailed: 4009, // Ошибка при создании или удалении индекса индекса в базе данных
}

type PanicCodeGroup struct {
	UnhandledException   int
	SecurityBreach       int
	SystemUnavailability int
}

var PANIC = PanicCodeGroup{
	UnhandledException:   5001, // Необработанное исключение
	SecurityBreach:       5002, // Критическая ошибка безопасности
	SystemUnavailability: 5003, // Полная недоступность системы
}

// Combinations for LogM (log multi)
var (
	MLvlAuditAll = []LogLevel{LvlAuditAPI, LvlAuditContent, LvlAuditPerms, LvlAuditCLI}
)

func ValidateEventCode(code int) error {
	validCodes := []int{
		INFO.AppStarted, INFO.UserLoggedIn, INFO.DbQuery, INFO.TaskCompleted,
		INFO.ConfigChanged, INFO.ServiceStarted, INFO.ServiceStopped, INFO.OperationCompleted,
		WARN.SuspiciousActivity, WARN.SlowQuery, WARN.LowResource, WARN.ParsingError,
		WARN.UnsupportedConfig, WARN.HighApiLatency,
		ERROR.DbConnectionFailed, ERROR.QueryExecutionFailed, ERROR.IOError,
		ERROR.AuthFailed, ERROR.DataIntegrationFailed,
		CRITICAL.CriticalServiceFailed, CRITICAL.DataLoss, CRITICAL.ImmediateAttention,
		CRITICAL.AppCrash, CRITICAL.BusinessProcessFailed,
		PANIC.UnhandledException, PANIC.SecurityBreach, PANIC.SystemUnavailability,
	}

	for _, validCode := range validCodes {
		if code == validCode {
			return nil
		}
	}

	return errors.New(fmt.Sprintf("Invalid event code: %d", code))
}
