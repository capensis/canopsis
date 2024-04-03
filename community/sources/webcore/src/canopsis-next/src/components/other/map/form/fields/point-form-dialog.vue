<template>
  <v-form
    v-click-outside.zIndex="showConfirmationModal"
    @submit.prevent.stop="submit"
  >
    <v-card width="400">
      <v-card-title class="primary pa-2 white--text">
        <v-layout
          justify-space-between
          align-center
        >
          <h4>{{ title }}</h4>
          <v-btn
            class="ma-0 ml-3"
            icon
            small
            @click="close"
          >
            <v-icon color="white">
              close
            </v-icon>
          </v-btn>
        </v-layout>
      </v-card-title>
      <v-card-text>
        <point-form
          v-model="form"
          :coordinates="coordinates"
          :exists-entities="existsEntities"
          @fly:coordinates="$emit('fly:coordinates', $event)"
        />
      </v-card-text>
      <v-card-actions class="v-layout justify-end">
        <v-btn
          :disabled="submitting"
          depressed
          text
          @click="close"
        >
          {{ $t('common.cancel') }}
        </v-btn>
        <v-btn
          v-if="editing"
          :disabled="submitting"
          :outlined="$system.dark"
          color="error"
          depressed
          text
          @click.stop="$emit('remove')"
        >
          {{ $t('common.delete') }}
        </v-btn>
        <v-btn
          :disabled="isDisabled"
          :loading="submitting"
          class="primary"
          type="submit"
        >
          {{ $t('common.submit') }}
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-form>
</template>

<script>
import { isEqual } from 'lodash';

import { MODALS } from '@/constants';

import { submittableMixinCreator } from '@/mixins/submittable';

import PointForm from './point-form.vue';

export default {
  $_veeValidate: {
    validator: 'new',
  },
  inject: ['$system'],
  components: { PointForm },
  mixins: [
    submittableMixinCreator(),
  ],
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
  data() {
    return {
      form: { ...this.point },
    };
  },
  computed: {
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

    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
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
