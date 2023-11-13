//About.vue
<template>
  <div class="about">
    <h1>This is an about page</h1>
    <button @click="getData()">test axios 请求数据</button>
    <p>这是请求到的二维码 </p>
    <img :src="testData.qrUrl">
  </div>
</template>
<script>
import { reactive } from 'vue'
import API from "../plugins/axiosInstance"
export default{
  name:'Qr',
  setup(){
    //数据
    const testData = reactive({
      qrUrl:"asfas"
    });

    //测试请求方法
    const getData = function(){
      API({
        url:'/qrUrl',
        method:'get'
      }).then((res)=>{
        testData.qrUrl = res.data.qrUrl;
        console.log(res.data.qrUrl)
      });
    }
    return{
      testData,
      getData,
    }
  }

}
</script>

