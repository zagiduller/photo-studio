<template>
  <div class="row mt-3">
    <div v-b-hover="imgHover" class="col-sm-12 col-md-6">
      <b-overlay :show="imgHovered" opacity="0.7">
        <b-img-lazy src="https://picsum.photos/1000/800" fluid alt="Responsive image" />
        <template #overlay>
          <div class="text-center">
            <b-icon icon="bell" font-scale="3" animation="cylon"></b-icon>
            <p id="cancel-label">если у вас уже есть заказ...</p>
            <b-button
              class="mt-5"
              variant="outline-success"
              lg="4"
              size="lg"
            >
              ВХОД
            </b-button>
          </div>
        </template>
      </b-overlay>
    </div>
    <div class="col-sm-12 col-md-6 mt-md-3">
      <div class="row justify-content-sm-center mt-2 ">
        <div class="col-sm-10 align-self-center">
          <div class="jumbotron">
            <h1 class="display-4">
              Новый заказ
            </h1>
            <div class="row">
              <div class="col-sm-12">
                <div class="form-group">
                  <label for="inputEmail">Почта</label>
                  <input id="exampleInputEmail1" type="email" class="form-control" aria-describedby="emailHelp">
                  <small id="emailHelp" class="form-text text-muted">Адрес электронной почты для связи.</small>
                </div>
                <div class="form-group">
                  <label for="inputPhone">Телефон</label>
                  <input id="inputPhone" type="text" class="form-control" aria-describedby="phoneHelp">
                  <small id="phoneHelp" class="form-text text-muted">Номер телефона для связи.</small>
                </div>
                <div class="form-group">
                  <label for="exampleFormControlTextarea1">Комментарии к заказу</label>
                  <textarea id="exampleFormControlTextarea1" class="form-control" rows="3" />
                </div>
              </div>
            </div>
            <hr class="my-4">
            <p>Опишите ваш заказ, оставьте почту или телефон и мы с вами свяжемся</p>
            <a class="btn btn-primary btn-lg" href="#" role="button">Отправить заявку</a>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'IndexPage',
  data () {
    return {
      city: 0,
      cities: [
        'Сибай',
        'Уфа',
        'Магнитогорск',
        'Белорецк'
      ],
      imgHovered: false
    }
  },
  computed: {
    currentCity () {
      return this.cities[this.city]
    }
  },
  methods: {
    imgHover (isHovered) {
      this.imgHovered = isHovered
    },
    prevCity () {
      this.city--
      if (this.city < 0) {
        this.city = this.cities.length - 1
      }
    },
    nextCity () {
      this.city++
      if (this.city >= this.cities.length) {
        this.city = 0
      }
    }
  }
}
</script>

<style>
</style>
