import { AxiosError } from 'axios'
import { ElMessage } from 'element-plus'

// APIのエラーレスポンスの構造を定義
interface ApiErrorResponse {
  message: string
}

export const showErrorMessage = (err: AxiosError) => {
  // unknown 型の data を ApiErrorResponse にキャスト
  const data = err.response?.data as ApiErrorResponse | undefined;
  const errorMessage = data?.message || '予備のエラーメッセージ';

  ElMessage({
    message: `エラーが発生しました\n${errorMessage}`,
    type: 'error'
  })
}