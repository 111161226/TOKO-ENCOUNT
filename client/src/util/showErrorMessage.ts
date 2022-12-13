import { AxiosError } from 'axios'
import { ElMessage } from 'element-plus'

export const showErrorMessage = (err: AxiosError) => {
  ElMessage({
    message: `エラーが発生しました\n${err.response?.data.message}`,
    type: 'error'
  })
}
