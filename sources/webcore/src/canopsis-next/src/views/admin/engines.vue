<template lang="pug">
  v-container
    h2.text-xs-center.my-3.display-1.font-weight-medium {{ $t('common.engines') }}
    engines-list(:loading="pending", :engines="engines")
    fab-buttons(@refresh="fetchList")
</template>

<script>
import entitiesEngineRunInfoMixin from '@/mixins/entities/engine-run-info';

import FabButtons from '@/components/other/fab-buttons/fab-buttons.vue';
import EnginesList from '@/components/other/engines/exploitation/engines-list.vue';

export default {
  components: { EnginesList, FabButtons },
  mixins: [entitiesEngineRunInfoMixin],
  data() {
    return {
      pending: true,
      engines: [],
    };
  },
  mounted() {
    this.fetchList();
  },
  methods: {
    async fetchList() {
      this.pending = true;

      this.engines = await this.fetchEnginesListWithoutStore();

      this.pending = false;
    },
  },
};
</script>

<style lang="scss" scoped>
.v-stepper {
  display: inline-block;
  box-shadow: none;
  background: none;

  &.event-info {
    float: right;
  }

  &__step__step {
    height: 40px;
    width: 40px;
    max-width: 40px;
    margin-right: 30px;
    font-size: 15px;
    font-weight: 500;
  }

  &.v-stepper--vertical &__content {
    margin: -8px -36px -16px 20px;
  }

  &:not(.v-stepper--vertical) &__label {
    display: flex !important;
  }

  .spacer {
    margin: 0 30px;
    width: 150px;
    border-top: 1px solid rgba(0, 0, 0, .12)
  }

  @media only screen and (max-width: 650px) {
    &.v-stepper--vertical &__content {
      margin: 0px -36px 0px 20px;
    }

    &:not(.v-stepper--vertical) &__content {
      padding-left: 0;
    }

    .v-stepper__step {
      padding: 10px;

      &__step {
        margin-right: 10px;
      }
    }

    .spacer {
      width: 50px;
    }
  }
}
</style>
