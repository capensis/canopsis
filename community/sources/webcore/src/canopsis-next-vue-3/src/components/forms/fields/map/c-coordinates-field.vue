<template>
  <c-layout :column="!row">
    <c-number-field
      :value="value.lat"
      :class="{ 'mr-3': row }"
      :label="$t('map.latitude')"
      :name="`${name}.lat`"
      :disabled="disabled"
      :min="-90"
      :max="90"
      :step="step"
      :required="isFieldRequired"
      @paste="pasteHandler('lat', $event)"
      @update:value="updateField('lat', $event)"
    />
    <c-number-field
      :value="value.lng"
      :label="$t('map.longitude')"
      :name="`${name}.lng`"
      :disabled="disabled"
      :min="-180"
      :max="180"
      :step="step"
      :required="isFieldRequired"
      @paste="pasteHandler('lng', $event)"
      @update:value="updateField('lng', $event)"
    />
  </c-layout>
</template>

<script>
import { isNumber, isString } from 'lodash';

import { formMixin } from '@/mixins/form';

export default {
  mixins: [formMixin],
  emits: ['input', 'update:value'],
  props: {
    value: {
      type: Object,
      default: () => ({
        lat: undefined,
        lng: undefined,
      }),
    },
    name: {
      type: String,
      default: 'coordinates',
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    step: {
      type: Number,
      default: 0.01,
    },
    row: {
      type: Boolean,
      default: false,
    },
    required: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    isFieldRequired() {
      return this.required || isNumber(this.value.lat) || isNumber(this.value.lng);
    },
  },
  methods: {
    pasteHandler(name, event) {
      const data = (event.clipboardData || window.clipboardData).getData('text');

      if (data && isString(data)) {
        if (data.includes(',')) {
          const [latString, lngString] = data.split(',');
          const lat = +latString;
          const lng = +lngString;

          if (isNumber(lat) && isNumber(lng)) {
            event.preventDefault();

            this.updateModel({ lat, lng });
          }
        }
      }
    },
  },
};
</script>
