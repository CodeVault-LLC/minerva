/* eslint-disable @typescript-eslint/naming-convention -- Enums */
export enum Category {
  SENSITIVE_DATA = "Sensitive Data",
  GAME_VARIABLE = "Game Variable",
  FUNCTION_NAME = "Function Name",
  GENERIC_NAME = "Generic Name",
  URL_PATH = "URL/Path",
  SECURITY_RISK = "Security Risk",
  VARIABLE_DECLARATION = "Variable Declaration",
  JSON_PROPERTY = "JSON/Property",
  JS_PATTERN = "JavaScript Pattern",
  COMMENT = "Comment",
  REGEX_LITERAL = "Regex Literal",
}

export enum SubCategory {
  API_KEY = "API Key",
  TOKEN = "Token",
  PASSWORD = "Password",
  GAME_STAT = "Game Statistic",
  GAME_ACTION = "Game Action",
  COMMON_VARIABLE = "Common Variable",
  URL = "URL",
  FILE_PATH = "File Path",
  DOM_MANIPULATION = "DOM Manipulation",
  CONSOLE_LOG = "Console Log",
  LOOP = "Loop",
  AJAX = "AJAX",
  BLOCK_COMMENT = "Block Comment",
  LINE_COMMENT = "Line Comment",
}
