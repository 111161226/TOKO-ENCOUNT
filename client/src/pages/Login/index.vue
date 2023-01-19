<script lang="ts" setup>
import { AxiosError } from 'axios'
import { reactive, ref, watchEffect } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, FormInstance } from 'element-plus'
import { useMe } from '@/store/me'
import { UserSimple } from '@/lib/apis'
import { showErrorMessage } from '@/util/showErrorMessage'
import { getRules } from '@/util/validate'

const meStore = useMe()

const formRef = ref<FormInstance>()
const rules = reactive(getRules(['userName', 'password']))
const isFormValid = ref(false)
watchEffect(() => {
  const { value } = formRef
  if (!value) {
    isFormValid.value = false
    return
  }

  if (inputData.password.length === 0) {
    isFormValid.value = false
    return
  }

  value.validate(isValid =>
    isValid ? (isFormValid.value = true) : (isFormValid.value = false)
  )
})

const inputData = reactive<UserSimple>({
  userName: '',
  password: ''
})
const loading = ref(false)
const router = useRouter()
const login = async () => {
  if (!isFormValid.value) {
    return
  }

  try {
    loading.value = true
    await meStore.login(inputData)
    ElMessage({
      message: 'ログインに成功しました',
      type: 'success'
    })
    router.push({ name: 'Home' })
  } catch (e: any) {
    const err: AxiosError = e
    showErrorMessage(err)
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="change-user-name-container">
    <div class="title">TOKO ENCOUNT</div>
    <el-form
      ref="formRef"
      :model="inputData"
      :rules="rules"
      label-position="top"
    >
      <el-form-item prop="userName" label="ユーザー名">
        <el-input v-model="inputData.userName" maxlength="30" show-word-limit/>
      </el-form-item>
      <el-form-item prop="password" label="パスワード">
        <el-input
          v-model="inputData.password"
          type="password"
          show-password
          @keyup.enter="login"
        />
      </el-form-item>
    </el-form>
    <div class="button-position">
    <el-button
      class="button"
      type="primary"
      :loading="loading"
      :disabled="!isFormValid"
      @click="login"
    >
      ログイン
    </el-button>
  </div>
    <div class="bottom-nav">
      <router-link :to="{ name: 'CreateAccount' }" class="link">アカウント作成</router-link>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.change-user-name-container {
  max-width: 600px;
  width: 80%;
  margin: 0 auto;
  padding-top: 20vh;
  

  .title {
    font-size: 30px;
    font-weight: bold;
    margin-bottom: 20px;
    text-align: center;
  }

  .bottom-nav {
    display: flex;
    justify-content: flex-end;
    margin-top: 30px;
  }

  .link{
    transition: color 0.2s;
    &:hover {
      color: $color-secondary;
    }
  }

  .button-position{
    display: flex;
    justify-content: center;
  }

  .button {
    background-color: $color-primary;
    color: white;
    border-radius: 0.5rem;
    width: 100%;
  }
}


</style>
