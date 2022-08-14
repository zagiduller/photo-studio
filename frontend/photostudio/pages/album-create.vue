<template>
  <div>
    <a-layout has-slider>
      <a-layout-sider theme="light">
        <a-menu mode="vertical">
          <a-menu-item key="1">Фото 1</a-menu-item>
          <a-menu-item key="2">Фото 2</a-menu-item>
          <a-menu-item key="3">Фото 3</a-menu-item>
        </a-menu>
      </a-layout-sider>
      <a-layout-content  class="a3" >
        <a-row type="flex" :gutter="[8,10]" style="padding: 5px">
          <template v-for="(block, id) in blocks">
            <a-col :key="id" span="6">
              <img width="390" :src="block.photo" />
            </a-col>
          </template>
        </a-row>
      </a-layout-content>
    </a-layout>
  </div>
</template>

<script>
export default {
  name: 'AlbumCreatePage',
  mounted () {
    const a5height = Math.ceil(595 / 2)
    const a5width = Math.ceil(842 / 2)

    // Todo: Надо переделать - разделить на две страницы

    for (let i = 0; i < 8; i++) {
      this.blocks.push({
        id: i,
        width: a5width,
        height: a5height,
        title: 'Блок ' + i,
        content: 'Содержимое блока ' + i,
        photo: 'https://picsum.photos/' + a5width + '/' + a5height + '?random'
      })
    }
  },
  data () {
    return {
      a4width: 842,
      a4height: 595,
      blocks: []
    }
  },
  methods: {
    calcWidth (count) {
      return this.a4width / count
    },
    calcHeight (count) {
      return this.a4height / count
    },
    clickBlock (event) {
      console.log(event)
    }
  },
  computed: {

  }
}
</script>

<style>
body {
  background-color: silver;
}
.a3 {
  width: 1700px;
  height: 1190px;
  background-color: cadetblue;
  /*padding: 5px 30px;*/
}
.a4 {
  height: 842px;
  width: 595px;
}
.a6 {
  height: 211px;
  width: 149px;
}
.block {
  background-color: white;
  /*margin: 2px;*/
}
</style>
