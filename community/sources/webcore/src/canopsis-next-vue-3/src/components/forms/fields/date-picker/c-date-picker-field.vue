<template>
  <v-menu
    v-model="opened"
    :close-on-content-click="false"
    :disabled="disabled"
    content-class="date-picker"
    transition="slide-y-transition"
    max-width="360"
    location="bottom end"
  >
    <template #activator="{ props: activatorProps }">
      <c-text-field
        v-bind="activatorProps"
        :class="contentClass"
        :model-value="textValue"
        :label="label"
        :placeholder="placeholder"
        :error="error"
        :name="name"
        :disabled="disabled"
        :hide-details="hideDetails"
        :append-inner-icon="clearable ? 'close' : ''"
        :readonly="!disabled"
        @click:append-inner="clear"
      >
        <template #append>
          <slot name="append" />
        </template>
      </c-text-field>
    </template>

    <v-date-picker
      :model-value="pickerValue"
      :color="color"
      :header-color="headerColor"
      :control-variant="controlVariant"
      :min="min"
      :max="max"
      :allowed-dates="allowedDates"
      :hide-header="hideHeader"
      class="date-picker"
      @update:model-value="updateDate"
    />
  </v-menu>
</template>

<script>
import { DATETIME_FORMATS } from '@/constants'
import { convertDateToDateObject, convertDateToString } from '@/helpers/date/date'
import { formBaseMixin } from '@/mixins/form'

export default {
  emits: ['update:value'],
  mixins: [formBaseMixin],
  props: {
    value: {
      type: [String, Number, Date],
      default: null,
    },
    label: {
      type: String,
      default: '',
    },
    placeholder: {
      type: String,
      default: '',
    },
    name: {
      type: String,
      default: '',
    },
    color: {
      type: String,
      default: 'primary',
    },
    headerColor: {
      type: String,
      default: 'primary',
    },
    controlVariant: {
      type: String,
      default: undefined,
    },
    format: {
      type: String,
      default: 'short',
    },
    error: {
      type: Boolean,
      default: false,
    },
    hideDetails: {
      type: Boolean,
      default: false,
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    min: {
      type: [String, Number, Date],
      required: false,
    },
    max: {
      type: [String, Number, Date],
      required: false,
    },
    allowedDates: {
      type: Function,
      required: false,
    },
    clearable: {
      type: Boolean,
      default: false,
    },
    contentClass: {
      type: [String, Object, Array],
      required: false,
    },
  },
  data() {
    return {
      opened: false,
    }
  },
  computed: {
    hideHeader() {
      return this.controlVariant !== 'modal'
    },

    textValue() {
      return convertDateToString(this.value, this.format)
    },

    pickerValue() {
      return this.value ? convertDateToDateObject(this.value, DATETIME_FORMATS.datePicker) : null
    },
  },
  methods: {
    clear() {
      this.updateModel(null)
    },

    updateDate(value) {
      const date = value instanceof Date ? value : convertDateToDateObject(value)
      const preparedValue = convertDateToString(date, DATETIME_FORMATS.datePicker, null)

      this.updateModel(preparedValue)
      this.$emit('change', preparedValue)
      this.opened = false
    },
  },
}
</script>

<style lang="scss">
.date-picker {
  .v-date-picker-month__day {
    height: 36px;
  }

  .v-date-picker-month__day .v-btn.v-date-picker-month__day-btn {
    --v-btn-height: 32px;
    --v-btn-size: 0.875rem;

    width: 32px;
    min-width: 32px;
    height: 32px;
  }

  .v-picker__body,
  .v-time-picker-clock__item,
  .v-time-picker-clock__item span,
  .v-time-picker-clock__hand {
    z-index: inherit;
  }
}
</style>
