<script lang="ts" setup>
import { AxiosError } from 'axios'
import { reactive, ref, watchEffect } from 'vue'
import { ElMessage, FormInstance } from 'element-plus'
import { useMe } from '@/store/me'
import { UserUpdate } from '@/lib/apis'
import { showErrorMessage } from '@/util/showErrorMessage'
import { getRules } from '@/util/validate'
import { prefectures } from '@/util/prefectures'
const meStore = useMe()
const formRef = ref<FormInstance>()
const rules = reactive(getRules(['userName','newPassword', 'password', 'prefect', 'gender']))
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

const inputData = reactive<UserUpdate>({
  userName: meStore.getMe?.userName ?? '',
  prefect: meStore.getMe?.prefect ?? '',
  gender: meStore.getMe?.gender ?? '',
  password: '',
  newPassword: ''
})

const loading = ref(false)
const confirmUpdate = async () => {
  if (!isFormValid.value) {
    return
  }

  try {
    loading.value = true
    await meStore.changeMeData({
      ...inputData,
      newPassword: inputData.password
    })
    ElMessage({
      message: 'ユーザー名を更新しました',
      type: 'success'
    })
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
    <div class="title">Profile</div>
    <el-form
      ref="formRef"
      :model="inputData"
      :rules="rules"
      label-position="top"
    >
      <el-form-item prop="userName" label="ユーザー名">
        <el-input v-model="inputData.userName" maxlength="30" show-word-limit/>
      </el-form-item>
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
      <el-form-item prop="newPassword" label="新しいパスワード">
        <el-input
            v-model="inputData.newPassword"
            type="password"
            show-password
            placeholder="変更しない場合は現在のパスワード"
            />
      </el-form-item>
      <el-form-item prop="password" label="現在のパスワードを入力して更新">
        <el-input
          v-model="inputData.password"
          type="password"
          show-password
          @keyup.enter="confirmUpdate"
        />
      </el-form-item>
    </el-form>

    <el-button
      type="primary"
      class="button"
      :loading="loading"
      :disabled="!isFormValid"
      @click="confirmUpdate"
    >
      更新
    </el-button>
  </div>
</template>

<style lang="scss" scoped>
.change-user-name-container {
  max-width: 600px;
  width: 80%;
  margin: 0 auto;
  padding: 40px 30px;

  .title {
    font-size: 24px;
    font-weight: bold;
    margin-bottom: 20px;
    text-align: center;
  }
}


.button {
  background-color: $color-primary;
  color: white;
  border-radius: 0.5rem;
  width: 100%;
  border: none;
  &:hover {
    background-color: $color-secondary;
  }
}
</style>