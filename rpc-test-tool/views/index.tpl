<!DOCTYPE html>
<html>
<head>
  <meta charset="utf-8" />
  <meta http-equiv="X-UA-Compatible" content="IE=edge">
  <title>配置</title>
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <!-- 引入element样式 -->
  <link rel="stylesheet" href="https://unpkg.com/element-ui/lib/theme-chalk/index.css">
  <style>
    #app {
      margin-top: 10px;
    }
    .el-input {
      margin-right: 10px;
    }
    .response {
      width: 600px;
      margin: 10px auto;
    }
  </style>
</head>
<body>
  <div id="app">
    <el-row>
      <el-col :span="10" :offset="7">
        <h2>配置</h2>
        <el-card>
          <el-form ref="form" :model="form" :rules="rules" size="small" label-position="right" label-width="80px">
            <el-form-item label="服务名称" prop="serveName">
              <el-input v-model="form.serveName" placeholder="服务名称"></el-input>
            </el-form-item>
            <el-form-item label="方法名称" prop="funcName">
              <el-input v-model="form.funcName" placeholder="方法名称"></el-input>
            </el-form-item>
            <div v-for="(arg, index) in form.args" :key="arg.key">
              <el-form-item
                style="display: inline-block"
                :label="'参数名'"
                :prop="'args.' + index + '.name'"
                :rules="{
                  required: true, message: '参数名不能为空', trigger: 'blur'
                }"
              >
                <el-input v-model="arg.name" placeholder="参数名"></el-input>
              </el-form-item>
              <el-form-item
                style="display: inline-block"
                :label="'参数值'"
                :key="arg.key"
                :prop="'args.' + index + '.value'"
                :rules="{
                  required: true, message: '参数不能为空', trigger: 'blur'
                }"
              >
                <el-input v-model="arg.value" placeholder="参数值"></el-input>
              </el-form-item>
              <el-button @click.prevent="removeArg(arg)" size="small">删除</el-button>
            </div>
            <el-form-item>
              <el-button type="primary" @click="submit">提交</el-button>
              <el-button @click="addArg">新增参数</el-button>
              <el-button @click="cancel">取消</el-button>
            </el-form-item>
          </el-form>
        </el-card>
      </el-col>
    </el-row>
    <div class="response" v-html="data"></div>
  </div>
  <!-- 引入vue -->
  <script src="https://cdn.jsdelivr.net/npm/vue"></script>
  <!-- 引入组件库 -->
  <script src="https://unpkg.com/element-ui/lib/index.js"></script>
  <!-- axios -->
  <script src="https://unpkg.com/axios/dist/axios.min.js"></script>
  <script>
    var app = new Vue({
      el: '#app',
      data: {
        form: {
          serveName: '',
          funcName: '',
          args: [{
            name: '',
            value: ''
          }]
        },
        rules: {
          serveName: [
            { required: true, message: '请输入服务名称', trigger: 'blur' }
          ],
          funcName: [
            { required: true, message: '请输入方法名称', trigger: 'blur' }
          ]
        },
        data: ''
      },
      methods: {
        submit() {
          this.$refs['form'].validate((valid) => {
          if (valid) {
            const data = {
              serveName: this.form.serveName,
              funcName: this.form.funcName,
              args: this.form.args.map(v => {
                return {
                  name: v.name,
                  value: v.value
                }
              })
            };
			const obj = this
            axios.post('http://www.rpctesttool.com:8080/start', data)
            .then(function (response) {
			  console.log(response);
              obj.cancel();
              obj.data = response.data;
            })
            .catch(function (error) {
              console.log(error);
              obj.data = error
            });
          } else {
            console.log('error submit!!');
            return false;
          }
        });
        },
        cancel() {
          this.$refs['form'].resetFields();
          this.form.args = [{
            name: '',
            value: ''
          }]
        },
        removeArg(item) {
          var index = this.form.args.indexOf(item)
          if (index !== -1) {
            this.form.args.splice(index, 1)
          }
        },
        addArg() {
          this.form.args.push({
            value: '',
            key: Date.now()
          });
        }
      }
    })
  </script>
</body>
</html>