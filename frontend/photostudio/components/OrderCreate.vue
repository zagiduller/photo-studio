<template>
  <div>
    <b-modal id="order-result" hide-footer>
      <template #modal-title>
        Результат
      </template>
      <div class="d-block text-center">
        <h3>{{ order }}</h3>
        <h3>{{ error }}</h3>
      </div>
      <b-button class="mt-3" block @click="$bvModal.hide('order-result')">
        Закрыть
      </b-button>
    </b-modal>
   <div class="col-sm-12 mt-5">
    <b-form-group
      id="inputName"
      label-cols-lg="2"
      content-cols-lg="10"
      label-cols-sm="12"
      content-cols-sm="12"
      description="Как к Вам обращаться?"
      label="Имя"
      label-for="input-horizontal"
      valid-feedback="Готово!"
    >
      <b-form-input id="inputName" v-model="name" type="text" />
    </b-form-group>
    <div class="row mt-5">
      <div class="col-md-6 col-sm-12">
        <b-form-group
          id="inputPhone"
          label-cols-sm="12"
          content-cols-sm="12"
          label-cols-lg="4"
          content-cols-lg="8"
          label="Телефон"
          label-for="input-horizontal"
          valid-feedback="Готово!"
          :invalid-feedback="invalidPhoneFeedback"
          :state="statePhone"
        >
          <b-form-input id="inputPhone" v-model="phone" size="sm" type="tel" />
        </b-form-group>
      </div>
      <div class="col-md-6 col-sm-12">
        <b-form-group
          id="inputEmail"
          label-cols-lg="4"
          content-cols-lg="8"
          label-cols-sm="12"
          content-cols-sm="12"
          label="Почта"
          label-for="input-horizontal"
          label-align-md="right"
          label-class="pr-md-4"
          valid-feedback="Готово!"
          :invalid-feedback="invalidEmailFeedback"
          :state="stateEmail"
        >
          <b-form-input id="inputEmail" v-model="email" size="sm" type="email" />
        </b-form-group>
      </div>
    </div>
     <hr>
     <b-form-group id="input-group-3" label="Город:" label-for="input-3">
       <b-form-select
         id="input-3"
         v-model="city"
         :options="cities"
         size="sm"
         required
       ></b-form-select>
     </b-form-group>
     <hr>
    <b-form-group
      id="descriptionFieldset"
      class="mt-5"
      label-for="inputDescription"
      label="Комментарий к заказу"
    >
      <b-form-textarea id="inputDescription" v-model="description" rows="5" type="tel" />
    </b-form-group>
  </div>
    <hr class="my-4">
    <p>Опишите ваш заказ, оставьте телефон или почту и мы с вами свяжемся</p>
    <b-button
      size="lg"
      :variant="!stateEmail && !statePhone ? '' : 'primary'"
      :disabled="!stateEmail && !statePhone"
      @click="newOrder"
    >
      Отправить
    </b-button>
  </div>
</template>

<script>
export default {
  name: 'OrderCreate',
  data () {
    return {
      order: {},
      error: '',
      name: '',
      email: '',
      phone: '',
      city: '',
      description: '',
      cities: [
        'Сибай',
        'Баймак',
        'Магнитогорск',
        'Белорецк',
        'Уфа'
      ]
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
            city: this.city,
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

<style scoped>

</style>
