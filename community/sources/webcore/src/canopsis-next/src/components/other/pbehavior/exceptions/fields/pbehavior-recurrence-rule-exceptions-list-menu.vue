<template>
  <v-menu offset-y>
    <template #activator="{ on }">
      <v-btn
        :loading="pending"
        :disabled="!availableExceptions.length"
        class="mr-0"
        color="primary"
        v-on="on"
      >
        {{ $t('pbehavior.exceptions.choose') }}
      </v-btn>
    </template>
    <v-list dense>
      <v-list-item
        v-for="exception in availableExceptions"
        :key="exception._id"
        @click="addItemIntoArray(exception)"
      >
        <v-list-item-content>
          <v-list-item-title>{{ exception.name }}</v-list-item-title>
        </v-list-item-content>
      </v-list-item>
    </v-list>
  </v-menu>
</template>

<script>
import { MAX_LIMIT } from '@/constants';

import { mapIds } from '@/helpers/array';

import { formArrayMixin } from '@/mixins/form';
import { entitiesPbehaviorExceptionMixin } from '@/mixins/entities/pbehavior/exceptions';

export default {
  inject: ['$validator'],
  mixins: [
    formArrayMixin,
    entitiesPbehaviorExceptionMixin,
  ],
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: Array,
      required: true,
    },
  },
  data() {
    return {
      pending: false,
      exceptions: [],
    };
  },
  computed: {
    selectedExceptionsIds() {
      return mapIds(this.value);
    },

    availableExceptions() {
      return this.exceptions.filter(({ _id: id }) => !this.selectedExceptionsIds.includes(id));
    },
  },
  mounted() {
    this.fetchList();
  },
  methods: {
    async fetchList() {
      this.pending = true;

      try {
        const { data } = await this.fetchPbehaviorExceptionsListWithoutStore({ params: { limit: MAX_LIMIT } });

        this.exceptions = data;
      } catch (err) {
        console.error(err);
      } finally {
        this.pending = false;
      }
    },
  },
};
</script>
