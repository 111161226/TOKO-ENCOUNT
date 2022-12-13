import validator from 'validator'
import { RuleItem } from 'async-validator'
import { FormRules } from 'element-plus'

export type InputField = 'password' | 'userName'

const isValidInput = (item: InputField, str: string) => {
  const validators: Record<InputField, (str: string) => boolean> = {
    password: str => validator.matches(str, /^[\x20-\x7E]{10,30}$/),
    userName: str => validator.matches(str, /^[a-zA-Z0-9_-]{1,30}$/)
  }

  return validators[item](str)
}

const errMessage = (item: InputField) => {
  const messages: Record<InputField, string> = {
    password: 'パスワードは10〜30文字で、半角英数字および半角記号が使えます',
    userName:
      'ユーザー名は1〜30文字で、英数字および"_" (半角アンダーバー), "-" (半角ハイフン) が使えます'
  }

  return messages[item]
}

const getValidator = (label: InputField): RuleItem['validator'] => {
  return (_rule, value, callback) => {
    if (!isValidInput(label, value)) {
      callback(errMessage(label))
    } else {
      callback()
    }
  }
}

export const getRules = (inputFields: InputField[]): FormRules =>
  Object.fromEntries(
    inputFields.map(item => [
      item,
      { validator: getValidator(item), trigger: 'blur' }
    ])
  )
