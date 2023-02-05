<template>
  <div class="row mt-3">
    <b-modal id="order-result" hide-footer>
      <template #modal-title>
        Результат
      </template>
      <div class="d-block text-center">
        <h3>{{ order }}</h3>
        <h3>{{ error }}</h3>
      </div>
      <b-button class="mt-3" block @click="$bvModal.hide('order-result')">
        Close Me
      </b-button>
    </b-modal>
    <div v-b-hover="imgHover" class="col-sm-12 col-md-6">
      <b-overlay :show="imgHovered" opacity="0.7">
        <b-img-lazy src="https://picsum.photos/1000/850" fluid alt="Здесь грузится картинка" />
        <template #overlay>
          <div class="text-center">
            <b-icon icon="bell" font-scale="3" animation="cylon" />
            <p id="cancel-label">
              если у вас уже есть заказ...
            </p>
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
    <div class="col-sm-12 col-md-6">
      <div class="row justify-content-sm-center mt-2 ">
        <div class="col-sm-10 align-self-center">
          <div v-if="haveUA">
            <div class="jumbotron">
              <b-button
                size="lg"
                :disabled="!haveUA"
                @click="haveUA = false"
              >
                У меня нет УЗ (учетный записи)
              </b-button>
            </div>
          </div>
          <div v-else>
            <OrderCreate />
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'IndexPage',
  components: {
    OrderCreate: () => import('~/components/OrderCreate.vue')
  },
  auth: false,
  data () {
    return {
      order: {},
      error: '',
      name: '',
      email: '',
      phone: '',
      description: '',
      imgHovered: false,
      haveUA: true
    }
  },
  computed: {
    stateEmail () {
      return this.email.length > 4
    },
    statePhone () {
      return this.phone.length > 4
    },
    invalidPhoneFeedback () {
      if (this.phone.length === 0) {
        if (this.stateEmail) {
          return ''
        }
        return 'Введите телефон для связи'
      }
      return 'Введите корректный телефон'
    },
    invalidEmailFeedback () {
      if (this.email.length === 0) {
        if (this.statePhone) {
          return ''
        }
        return 'Введите почту для связи'
      }
      return 'Почта должна быть правильной'
    }
  },
  methods: {
    imgHover (isHovered) {
      this.imgHovered = isHovered
    },
    async newOrder () {
      this.error = ''
      if (this.statePhone || this.stateEmail) {
        try {
          const order = await this.$axios.post('/orders/', {
            name: this.name,
            email: this.email,
            phone: this.phone,
            description: this.description
          })
          console.log(order)
          this.order = order
        } catch (e) {
          this.error = e.response.data
          console.log(e)
        }
        this.$bvModal.show('order-result')
      }
    }
  }
}
</script>

<style>
</style>
