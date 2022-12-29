<script lang="ts" setup>
import { AxiosError } from 'axios'
import { reactive, ref, watchEffect } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, FormInstance } from 'element-plus'
import { useMe } from '@/store/me'
import { User } from '@/lib/apis'
import { showErrorMessage } from '@/util/showErrorMessage'
import { getRules } from '@/util/validate'
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
      label-position="top"
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
          <el-option value="北海道" label="北海道" />
          <el-option value="青森県" label="青森県" />
          <el-option value="岩手県" label="岩手県" />
          <el-option value="宮城県" label="宮城県" />
          <el-option value="秋田県" label="秋田県" />
          <el-option value="山形県" label="山形県" />
          <el-option value="福島県" label="福島県" />
          <el-option value="茨城県" label="茨城県" />
          <el-option value="栃木県" label="栃木県" />
          <el-option value="群馬県" label="群馬県" />
          <el-option value="埼玉県" label="埼玉県" />
          <el-option value="千葉県" label="千葉県" />
          <el-option value="東京都" label="東京都" />
          <el-option value="神奈川県" label="神奈川県" />
          <el-option value="新潟県" label="新潟県" />
          <el-option value="富山県" label="富山県" />
          <el-option value="石川県" label="石川県" />
          <el-option value="福井県" label="福井県" />
          <el-option value="山梨県" label="山梨県" />
          <el-option value="長野県" label="長野県" />
          <el-option value="岐阜県" label="岐阜県" />
          <el-option value="静岡県" label="静岡県" />
          <el-option value="愛知県" label="愛知県" />
          <el-option value="三重県" label="三重県" />
          <el-option value="滋賀県" label="滋賀県" />
          <el-option value="京都府" label="京都府" />
          <el-option value="大阪府" label="大阪府" />
          <el-option value="兵庫県" label="兵庫県" />
          <el-option value="奈良県" label="奈良県" />
          <el-option value="和歌山県" label="和歌山県" />
          <el-option value="鳥取県" label="鳥取県" />
          <el-option value="島根県" label="島根県" />
          <el-option value="岡山県" label="岡山県" />
          <el-option value="広島県" label="広島県" />
          <el-option value="山口県" label="山口県" />
          <el-option value="徳島県" label="徳島県" />
          <el-option value="香川県" label="香川県" />
          <el-option value="愛媛県" label="愛媛県" />
          <el-option value="高知県" label="高知県" />
          <el-option value="福岡県" label="福岡県" />
          <el-option value="佐賀県" label="佐賀県" />
          <el-option value="長崎県" label="長崎県" />
          <el-option value="熊本県" label="熊本県" />
          <el-option value="大分県" label="大分県" />
          <el-option value="宮崎県" label="宮崎県" />
          <el-option value="鹿児島県" label="鹿児島県" />
          <el-option value="沖縄県" label="沖縄県" />
        </el-select>
      </el-form-item>
    </el-form>

    <el-button
      type="primary"
      :loading="loading"
      :disabled="!isFormValid"
      @click="confirmCreate"
    >
      作成
    </el-button>
    <div class="bottom-nav">
      <router-link :to="{ name: 'Login' }">ログイン</router-link>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.create-user-container {
  padding: 40px 30px;

  .title {
    font-size: 24px;
    font-weight: bold;
    margin-bottom: 20px;
    text-align: center;
  }

}
</style>
