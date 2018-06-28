<template lang="pug">
  v-container.pa-0(fluid)
    v-layout(wrap)
      v-flex.my-2(xs12)
        v-layout(justify-space-around, align-center, wrap)
          v-flex(xs4, md2)
            h4.text-xs-center {{ $t(`context.moreInfos.type`) }} :
            p.text-xs-center {{ item.props.type }}
          v-flex(xs4, md2)
            template(v-if="item.props.enabled")
              v-chip.green.white--text {{ $t(`common.enabled`) }}
            template(v-else)
              v-chip.red.white--text.title {{ $t(`common.disabled`) }}
          v-flex(xs4, md2)
            h4.text-xs-center {{ $t(`context.moreInfos.lastActiveDate`) }} :
            p.text-xs-center {{ this.$d(new Date( item.props.enable_history[0] * 1000), 'short') }}
          v-flex(xs6, md2)
            v-menu(:value="isImpactExpanded", bottom, offset-y, fixed)
              v-btn(@click.stop="isImpactExpanded = !isImpactExpanded", slot="activator") {{ $t(`context.impacts`) }}
              v-list(dense)
                template(v-for="item in item.props.impact")
                  v-list-tile
                    v-list-tile-content {{ item }}
          v-flex(xs6, md2)
            v-menu(:value="isDependsExpanded", bottom, offset-y, fixed)
              v-btn(@click.stop="isDependsExpanded = !isDependsExpanded",
                    slot="activator") {{ $t(`context.dependencies`) }}
              v-list(dense)
                template(v-for="item in item.props.depends")
                  v-list-tile
                    v-list-tile-content {{ item }}
      v-flex.my-2(xs12)
        h3.text-xs-center.my-2 Pbehaviors
        pbehaviors-list(:itemId="item.props._id")
      v-flex.my-2(xs12)
        h3.text-xs-center Infos
        v-container(fluid, grid-list-sm)
          v-layout(row, wrap)
            v-flex(v-for="(value, key) in lol.lol" xs4)
              h4.text-xs-center {{ key }}
              p.text-xs-center {{ $t(`common.description`) }} : {{ value.description }}
              p.text-xs-center {{ $t(`common.value`) }} : {{ value.value }}


</template>

<script>
import PbehaviorsList from './pbehaviors-list.vue';

export default {
  name: 'context-more-infos',
  components: {
    PbehaviorsList,
  },
  props: {
    item: {
      type: Object,
      required: true,
    },
  },
  data() {
    return {
      isImpactExpanded: false,
      isDependsExpanded: false,
      lol: {
        lol:
    {
      service_period: {
        description: "Plage de service de l'application",
        value: '',
      },
      manual_maintenance_comment: {
        description: "Commentaire ajouté lors d'une mise en maintenance manuelle",
        value: '',
      },
      display_on_weather: {
        description: "Afficher ou non l'application sur une météo",
        value: 'False',
      },
      application_crit_code: {
        description: 'Criticité - Code',
        value: '3',
      },
      weather_type: {
        description: 'Type de météo de service',
        value: 'MDSA',
      },
      application_label: {
        description: 'Libellé application',
        value: 'Gestion budgétaire du département pilotage GA SI',
      },
      manual_maintenance: {
        description: "Indique si l'app est actuellement en maintenance manuelle",
        value: 'False',
      },
    },
      },
    };
  },
  computed: {
    infos() {
      return this.item.props.infos || '';
    },
  },
};
</script>

