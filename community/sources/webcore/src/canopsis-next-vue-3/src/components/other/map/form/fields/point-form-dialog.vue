<template>
  <v-form v-click-outside.zIndex="showConfirmationModal" @submit.prevent.stop="submit">
    <v-card width="400">
      <v-card-title class="bg-primary text-white pa-2">
        <c-layout align-center justify-space-between>
          <h4>{{ title }}</h4>
          <v-btn class="ma-0 ml-3" size="small" variant="text" icon @click="close">
            <v-icon color="white"> close </v-icon>
          </v-btn>
        </c-layout>
      </v-card-title>
      <v-card-text>
        <point-form
          :form="form"
          :coordinates="coordinates"
          :exists-entities="existsEntities"
          @update:form="form = $event"
          @fly:coordinates="$emit('fly:coordinates', $event)"
        />
      </v-card-text>
      <v-card-actions class="justify-end">
        <v-btn :disabled="submitting" variant="text" @click="close">
          {{ $t('common.cancel') }}
        </v-btn>
        <v-btn
          v-if="editing"
          :disabled="submitting"
          color="error"
          variant="text"
          @click.stop="$emit('remove')"
        >
          {{ $t('common.delete') }}
        </v-btn>
        <v-btn :disabled="isDisabled" :loading="submitting" color="primary" type="submit">
          {{ $t('common.submit') }}
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-form>
</template>

<script>
import { isEqual } from 'lodash';

import { MODALS } from '@/constants';

import PointForm from './point-form.vue';

export default {
  components: { PointForm },
  props: {
    point: {
      type: Object,
      required: true,
    },
    editing: {
      type: Boolean,
      required: false,
    },
    coordinates: {
      type: Boolean,
      required: false,
    },
    existsEntities: {
      type: Array,
      default: () => [],
    },
  },
  emits: ['cancel', 'submit', 'remove', 'fly:coordinates'],
  data() {
    return {
      form: { ...this.point },
      submitting: false,
    };
  },
  computed: {
    isDisabled() {
      return !this.form.entity;
    },

    title() {
      if (this.editing) {
        return this.$t('map.editPoint');
      }

      return this.$t('map.addPoint');
    },
  },
  watch: {
    point(point) {
      this.form = { ...point };
    },
  },
  methods: {
    close() {
      this.$emit('cancel');
    },

    submit() {
      if (!this.isDisabled) {
        this.$emit('submit', this.form);
      }
    },

    showConfirmationModal() {
      if (!isEqual(this.point, this.form)) {
        this.$modals.show({
          id: this.point._id,
          name: MODALS.clickOutsideConfirmation,
          dialogProps: {
            persistent: true,
          },
          config: {
            action: (confirmed) => {
              if (confirmed) {
                return this.submit();
              }

              return this.close();
            },
          },
        });
      } else {
        this.close();
      }
    },
  },
};
</script>
