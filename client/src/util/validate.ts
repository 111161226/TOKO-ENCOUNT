import validator from 'validator'
import { RuleItem } from 'async-validator'
import { FormRules } from 'element-plus'

export type InputField = 'password' | 'newPassword'| 'userName' | 'prefect' | 'gender'

const isValidInput = (item: InputField, str: string) => {
  const validators: Record<InputField, (str: string) => boolean> = {
    password: str => validator.matches(str, /^[\x20-\x7E]{10,30}$/),
    newPassword: str => validator.matches(str, /^[\x20-\x7E]{10,30}$/),
    userName: str => validator.matches(str, /^[a-zA-Z0-9_-]{1,30}$/),
    prefect: str => validator.matches(str, /^[一-龠]{1,30}$/),
    gender: str => validator.matches(str, /^[a-zA-Z0-9_-]{1,30}$/)
  }

  return validators[item](str)
}

const errMessage = (item: InputField) => {
  const messages: Record<InputField, string> = {
    password: 'パスワードは10〜30文字で、半角英数字および半角記号が使えます',
    newPassword: 'パスワードは10〜30文字で、半角英数字および半角記号が使えます',
    userName:
      'ユーザー名は1〜30文字で、英数字および"_" (半角アンダーバー), "-" (半角ハイフン) が使えます',
    prefect: '都道府県を選んでください',
    gender: '性別を選んでください'
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
