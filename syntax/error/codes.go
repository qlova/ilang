package error

import "github.com/qlova/uct/compiler"

var Codes = [512]compiler.Translatable{
	1: compiler.Translatable{
		compiler.English: "Unknown",
	},
	
	2: compiler.Translatable{
		compiler.English: "DoesNotExist",
	},
	
	3: compiler.Translatable{
		compiler.English: "NotSupported",
	},
	
	4: compiler.Translatable{
		compiler.English: "Interrupted",
	},
	
	5: compiler.Translatable{
		compiler.English: "InputOutput",
	},
	
	6: compiler.Translatable{
		compiler.English: "Impossible",
	},
	
	7: compiler.Translatable{
		compiler.English: "TooManyArguments",
	},
	
	8: compiler.Translatable{
		compiler.English: "Incompatible",
	},
	
	9: compiler.Translatable{
		compiler.English: "BadNumber",
	},
	
	10: compiler.Translatable{
		compiler.English: "",
	},
	
	11: compiler.Translatable{
		compiler.English: "TryAgain",
	},
	
	12: compiler.Translatable{
		compiler.English: "OutOfMemory",
	},
	
	13: compiler.Translatable{
		compiler.English: "PermissionDenied",
	},
	
	14: compiler.Translatable{
		compiler.English: "BadAddress",
	},
	
	15: compiler.Translatable{
		compiler.English: "NoMoreData",
	},
	
	16: compiler.Translatable{
		compiler.English: "Busy",
	},
	
	17: compiler.Translatable{
		compiler.English: "Exists",
	},
	
	18: compiler.Translatable{
		compiler.English: "",
	},
	
	19: compiler.Translatable{
		compiler.English: "",
	},
	
	20: compiler.Translatable{
		compiler.English: "Directory",
	},
	
	21: compiler.Translatable{
		compiler.English: "NotReady",
	},
	
	22: compiler.Translatable{
		compiler.English: "InvalidArgument",
	},
	
	23: compiler.Translatable{
		compiler.English: "TooManyLinks",
	},
	
	24: compiler.Translatable{
		compiler.English: "TooManyOpen",
	},
	
	25: compiler.Translatable{
		compiler.English: "CannotSend",
	},
	
	26: compiler.Translatable{
		compiler.English: "CannotRead",
	},
	
	27: compiler.Translatable{
		compiler.English: "TooBig",
	},
	
	28: compiler.Translatable{
		compiler.English: "NoSpaceLeft",
	},
	
	29: compiler.Translatable{
		compiler.English: "Illegal",
	},
	
	30: compiler.Translatable{
		compiler.English: "ReadOnly",
	},
	
	31: compiler.Translatable{
		compiler.English: "TooBig",
	},
	
	32: compiler.Translatable{
		compiler.English: "Broken",
	},
	
	33: compiler.Translatable{
		compiler.English: "OutOfDomain",
	},
	
	34: compiler.Translatable{
		compiler.English: "NotRepresentable",
	},
	
	35: compiler.Translatable{
		compiler.English: "Deadlock",
	},
	
	36: compiler.Translatable{
		compiler.English: "NameTooLong",
	},
	
	37: compiler.Translatable{
		compiler.English: "NotAvailable",
	},
	
	38: compiler.Translatable{
		compiler.English: "",
	},
	
	39: compiler.Translatable{
		compiler.English: "NotEmpty",
	},
	
	40: compiler.Translatable{
		compiler.English: "",
	},
	
	41: compiler.Translatable{
		compiler.English: "",
	},
	
	42: compiler.Translatable{
		compiler.English: "CannotCommunicate",
	},
	
	43: compiler.Translatable{
		compiler.English: "Corrupted",
	},
	
	44: compiler.Translatable{
		compiler.English: "OutOfRange",
	},
	
	45: compiler.Translatable{
		compiler.English: "",
	},
	
	46: compiler.Translatable{
		compiler.English: "",
	},
	
	47: compiler.Translatable{
		compiler.English: "",
	},
	
	48: compiler.Translatable{
		compiler.English: "",
	},
	
	49: compiler.Translatable{
		compiler.English: "NotAttached",
	},
	
	50: compiler.Translatable{
		compiler.English: "",
	},
	
	51: compiler.Translatable{
		compiler.English: "",
	},
	
	52: compiler.Translatable{
		compiler.English: "LowEntropy",
	},
	
	53: compiler.Translatable{
		compiler.English: "",
	},
	
	54: compiler.Translatable{
		compiler.English: "",
	},
	
	55: compiler.Translatable{
		compiler.English: "NotEnoughPower",
	},
	
	56: compiler.Translatable{
		compiler.English: "",
	},
	
	57: compiler.Translatable{
		compiler.English: "",
	},
	
	59: compiler.Translatable{
		compiler.English: "BadFormat",
	},
	
	60: compiler.Translatable{
		compiler.English: "",
	},
	
	61: compiler.Translatable{
		compiler.English: "Empty",
	},
	
	62: compiler.Translatable{
		compiler.English: "OutOfTime",
	},

	63: compiler.Translatable{
		compiler.English: "",
	},
	
	64: compiler.Translatable{
		compiler.English: "",
	},
	
	65: compiler.Translatable{
		compiler.English: "NotInstalled",
	},
	
	66: compiler.Translatable{
		compiler.English: "TooFarAway",
	},
	
	67: compiler.Translatable{
		compiler.English: "",
	},
	
	68: compiler.Translatable{
		compiler.English: "NotUnique",
	},
	
	69: compiler.Translatable{
		compiler.English: "",
	},
	
	70: compiler.Translatable{
		compiler.English: "",
	},
	
	71: compiler.Translatable{
		compiler.English: "Protocol",
	},
	
	72: compiler.Translatable{
		compiler.English: "",
	},
	
	73: compiler.Translatable{
		compiler.English: "",
	},
	
	74: compiler.Translatable{
		compiler.English: "",
	},
	
	75: compiler.Translatable{
		compiler.English: "",
	},
	
	76: compiler.Translatable{
		compiler.English: "",
	},
	
	77: compiler.Translatable{
		compiler.English: "BadState",
	},
	
	78: compiler.Translatable{
		compiler.English: "",
	},
	
	79: compiler.Translatable{
		compiler.English: "MissingLibrary",
	},
	
	81: compiler.Translatable{
		compiler.English: "",
	},
	
	82: compiler.Translatable{
		compiler.English: "",
	},
	
	83: compiler.Translatable{
		compiler.English: "",
	},
	
	84: compiler.Translatable{
		compiler.English: "",
	},
	
	85: compiler.Translatable{
		compiler.English: "",
	},
	
	86: compiler.Translatable{
		compiler.English: "",
	},
	
	87: compiler.Translatable{
		compiler.English: "TooManyUsers",
	},
	
	88: compiler.Translatable{
		compiler.English: "WrongType",
	},
	
	89: compiler.Translatable{
		compiler.English: "MissingData",
	},
	
	90: compiler.Translatable{
		compiler.English: "",
	},
	
	91: compiler.Translatable{
		compiler.English: "",
	},
	
	92: compiler.Translatable{
		compiler.English: "",
	},
	
	93: compiler.Translatable{
		compiler.English: "",
	},
	
	94: compiler.Translatable{
		compiler.English: "",
	},
	
	95: compiler.Translatable{
		compiler.English: "",
	},
	
	96: compiler.Translatable{
		compiler.English: "",
	},
	
	97: compiler.Translatable{
		compiler.English: "",
	},
	
	98: compiler.Translatable{
		compiler.English: "Taken",
	},
	
	99: compiler.Translatable{
		compiler.English: "",
	},
	
	100: compiler.Translatable{
		compiler.English: "NetworkDown",
	},
	
	101: compiler.Translatable{
		compiler.English: "NetworkUnreachable",
	},
	
	102: compiler.Translatable{
		compiler.English: "NetworkDropped",
	},
	
	103: compiler.Translatable{
		compiler.English: "Abort",
	},
	
	104: compiler.Translatable{
		compiler.English: "Reset",
	},
	
	105: compiler.Translatable{
		compiler.English: "BufferFull",
	},
	
	106: compiler.Translatable{
		compiler.English: "AlreadyConnected",
	},
	
	107: compiler.Translatable{
		compiler.English: "NotConnected",
	},
	
	108: compiler.Translatable{
		compiler.English: "Shutdown",
	},
	
	110: compiler.Translatable{
		compiler.English: "Timeout",
	},
	
	111: compiler.Translatable{
		compiler.English: "ConnectionRefused",
	},
	
	112: compiler.Translatable{
		compiler.English: "HostDown",
	},
	
	113: compiler.Translatable{
		compiler.English: "MissingRoute",
	},
	
	114: compiler.Translatable{
		compiler.English: "",
	},
	
	115: compiler.Translatable{
		compiler.English: "",
	},
	
	116: compiler.Translatable{
		compiler.English: "",
	},
	
	117: compiler.Translatable{
		compiler.English: "NeedsCleaning",
	},
	
	118: compiler.Translatable{
		compiler.English: "",
	},
	
	119: compiler.Translatable{
		compiler.English: "",
	},
	
	120: compiler.Translatable{
		compiler.English: "",
	},
	
	121: compiler.Translatable{
		compiler.English: "",
	},
	
	122: compiler.Translatable{
		compiler.English: "QuotaExceeded",
	},
	
	123: compiler.Translatable{
		compiler.English: "",
	},
	
	124: compiler.Translatable{
		compiler.English: "",
	},
	
	400: compiler.Translatable{
		compiler.English: "BadRequest",
	},
	
	401: compiler.Translatable{
		compiler.English: "Unauthroized",
	},
	
	402: compiler.Translatable{
		compiler.English: "PaymentRequired",
	},
	
	403: compiler.Translatable{
		compiler.English: "Forbidden",
	},
	
	404: compiler.Translatable{
		compiler.English: "NotFound",
	},
	
	405: compiler.Translatable{
		compiler.English: "MethodNotAllowed",
	},
	
	406: compiler.Translatable{
		compiler.English: "NotAcceptable",
	},
	
	407: compiler.Translatable{
		compiler.English: "ProxyAuthenticationRequired",
	},
	
	408: compiler.Translatable{
		compiler.English: "RequestTimeout",
	},
	
	409: compiler.Translatable{
		compiler.English: "Conflict",
	},
	
	410: compiler.Translatable{
		compiler.English: "Gone",
	},
	
	411: compiler.Translatable{
		compiler.English: "LengthRequired",
	},
	
	412: compiler.Translatable{
		compiler.English: "PreconditionFailed",
	},
	
	413: compiler.Translatable{
		compiler.English: "PayloadTooLarge",
	},
	
	414: compiler.Translatable{
		compiler.English: "URITooLong",
	},
	
	415: compiler.Translatable{
		compiler.English: "UnsupportedMediaType",
	},
	
	416: compiler.Translatable{
		compiler.English: "RangeNotSatisfiable",
	},

	417: compiler.Translatable{
		compiler.English: "ExpectationFailed",
	},
	
	418: compiler.Translatable{
		compiler.English: "ImATeapot",
	},
	
	421: compiler.Translatable{
		compiler.English: "MisdirectedRequest",
	},
	
	422: compiler.Translatable{
		compiler.English: "UnprocessableEntity",
	},
	
	423: compiler.Translatable{
		compiler.English: "Locked",
	},
	
	424: compiler.Translatable{
		compiler.English: "FailedDependency",
	},
	
	426: compiler.Translatable{
		compiler.English: "UpgradeRequired",
	},
	
	428: compiler.Translatable{
		compiler.English: "PreconditionRequired",
	},
	
	429: compiler.Translatable{
		compiler.English: "TooManyRequests",
	},
	
	431: compiler.Translatable{
		compiler.English: "RequestHeaderFieldsTooLarge",
	},
	
	451: compiler.Translatable{
		compiler.English: "UnavailableForLegalReasons",
	},
	
	500: compiler.Translatable{
		compiler.English: "InternalServerError",
	},
	
	501: compiler.Translatable{
		compiler.English: "NotImplemented",
	},
	
	502: compiler.Translatable{
		compiler.English: "BadGateway",
	},
	
	503: compiler.Translatable{
		compiler.English: "ServiceUnavailable",
	},
	
	504: compiler.Translatable{
		compiler.English: "GatewayTimeout",
	},
	
	505: compiler.Translatable{
		compiler.English: "HTTPVersionNotSupported",
	},
	
	506: compiler.Translatable{
		compiler.English: "VariantAlsoNegotiates",
	},
	
	507: compiler.Translatable{
		compiler.English: "InsufficientStorage",
	},
	
	508: compiler.Translatable{
		compiler.English: "LoopDetected",
	},
	
	510: compiler.Translatable{
		compiler.English: "NotExtended",
	},
	
	511: compiler.Translatable{
		compiler.English: "NetworkAuthenticationRequired",
	},
}
