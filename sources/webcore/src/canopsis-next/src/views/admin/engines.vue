<template lang="pug">
  v-container
    h2.text-xs-center.my-3.display-1.font-weight-medium {{ $t('common.engines') }}
    v-fade-transition
      progress-overlay(v-if="pending", :pending="pending")
      v-layout(v-else, row, justify-center)
        v-flex
          div.v-stepper.theme--light.event-info
            div.v-stepper__step.v-stepper__step--active.pr-0
              span.v-stepper__step__step.blue.darken-4.ma-0
                v-icon(medium) offline_bolt
              div.spacer
            div.v-stepper__content.v-stepper__step--active.pt-0
              div.v-stepper__label
                div
                  strong {{ $t('engines.event.title') }}
                small {{ $t('engines.event.description') }}
        v-flex
          div.v-stepper.v-stepper--vertical.theme--light
            template(v-for="(engine, index) in engines")
              div.v-stepper__step.v-stepper__step--active.pl-0(:key="engine.name")
                span.v-stepper__step__step.primary {{ index + 1 }}
                div.v-stepper__label
                  div
                    strong {{ $t(`engines.${engine.name}.title`) }}
                  small {{ $t(`engines.${engine.name}.description`) }}
              div.v-stepper__content(:key="`${engine.name}-content`")
    fab-buttons(@refresh="fetchList")
</template>

<script>
import entitiesEngineRunInfoMixin from '@/mixins/entities/engine-run-info';

import ProgressOverlay from '@/components/layout/progress/progress-overlay.vue';
import FabButtons from '@/components/other/fab-buttons/fab-buttons.vue';

export default {
  components: { ProgressOverlay, FabButtons },
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
