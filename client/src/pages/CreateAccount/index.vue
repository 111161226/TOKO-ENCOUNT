<script lang="ts" setup>
import { AxiosError } from 'axios'
import { reactive, ref, watchEffect } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, FormInstance } from 'element-plus'
import { useMe } from '@/store/me'
import { User } from '@/lib/apis'
import { showErrorMessage } from '@/util/showErrorMessage'
import { getRules } from '@/util/validate'
import { prefectures } from '@/util/prefectures'
const meStore = useMe()

const formRef = ref<FormInstance>()
const rules = reactive(getRules(['userName', 'password','gender','prefect']))
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

const inputData = reactive<User>({
  userName: '',
  password: '',
  prefect: '東京都',
  gender: 'male'
})

const loading = ref(false)
const router = useRouter()
const confirmCreate = async () => {
  if (!isFormValid.value) {
    return
  }

  try {
    loading.value = true
    await meStore.createMe(inputData)
    ElMessage({
      message: 'アカウントを作成しました',
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
  <div class="create-user-container">
    <div class="title">アカウント作成</div>
    <el-form
      ref="formRef"
      :model="inputData"
      :rules="rules"
      label-position="left"
    >
      <el-form-item prop="userName" label="ユーザー名">
        <el-input v-model="inputData.userName" maxlength="30" show-word-limit />
      </el-form-item>
      <el-form-item prop="password" label="パスワード">
        <el-input
          v-model="inputData.password"
          type="password"
          maxlength="30"
          show-word-limit
          show-password
          @keyup.enter="confirmCreate"
        />
      </el-form-item>
      <div style="display:inline-flex">
      <el-form-item prop="gender" label="性別">
        <el-select
          v-model="inputData.gender"
        >
          <el-option label="male" value="male" />
          <el-option label="female" value="female" />
        </el-select>
      </el-form-item>
      <el-form-item prop="prefect" label="都道府県">
        <el-select v-model="inputData.prefect" > 
          <el-option
            v-for="item in prefectures"
            :key="item.value"
            :label="item.label"
            :value="item.value"
          />
        </el-select>
      </el-form-item>
      </div>
    </el-form>
    <div class="button">
    <el-button
      type="primary"
      :loading="loading"
      :disabled="!isFormValid"
      @click="confirmCreate"
    >
      作成
    </el-button>
  </div>
    <div class="bottom-nav">
      <router-link :to="{ name: 'Login' }">ログイン</router-link>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.create-user-container {
  width: 35%;
  margin: 0 auto;
  padding: 40px 30px;

  .title {
    font-size: 24px;
    font-weight: bold;
    margin-bottom: 20px;
    text-align: center;
  }
  
}
</style>
