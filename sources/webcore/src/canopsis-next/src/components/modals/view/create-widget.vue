<template lang='pug'>
  v-card(data-test="createWidgetModal")
    v-card-title.primary.white--text
      v-layout(justify-space-between, align-center)
        span.headline {{ $t('modals.widgetCreation.title') }}
        v-btn(icon, @click="hideModal")
          v-icon.white--text clear
    v-card-text.pa-0
      v-layout
        v-flex(xs6)
          v-list(dense)
            v-subheader.body-1.mb-1 Général
            div(v-for="widget in generalTypeWidgets", :key="widget.value")
              v-divider
              v-list-tile(@click.stop="selectWidget(widget.value)")
                v-list-tile-avatar
                  v-icon {{ widget.icon }}
                v-list-tile-content
                  v-list-tile-title {{ widget.title }}
                v-list-tile-action
                  v-btn.primary--text(icon, @click.stop="selectWidgetType(widget.value)")
                    v-icon add
          v-list(dense)
            v-subheader Stats
            v-divider
            div(v-for="widget in statTypeWidgets", :key="widget.value")
              v-divider
              v-list-tile(@click.stop="selectWidget(widget.value)")
                v-list-tile-avatar
                  v-icon {{ widget.icon }}
                v-list-tile-content
                  v-list-tile-title {{ widget.title }}
                v-list-tile-action
                  v-btn.primary--text(icon, @click.stop="selectWidgetType(widget.value)")
                    v-icon add
        v-flex(xs6)
          v-subheader.body-1 Aperçu
          v-divider
          v-container
            h4 {{ selectedWidget.title }}
            p.text-xs-justify {{ selectedWidget.description }}
            v-layout.px-4(wrap, v-if="selectedWidget")
              v-flex(xs3)
                v-img(
                  :src="`https://picsum.photos/500/300?image=1`",
                  :lazy-src="`https://picsum.photos/10/6?image=1`",
                  aspect-ratio="1"
                )
              v-spacer
              v-flex(xs3)
                v-img(
                  :src="`https://picsum.photos/500/300?image=2`",
                  :lazy-src="`https://picsum.photos/10/6?image=2`",
                  aspect-ratio="1"
                )
              v-spacer
              v-flex(xs3)
                v-img(
                  :src="`https://picsum.photos/500/300?image=3`",
                  :lazy-src="`https://picsum.photos/10/6?image=3`",
                  aspect-ratio="1"
                )
            div(v-else)
              v-layout(align-content-center, justify-center)
                span
                  v-icon(small) info
                p.ma-0 Please select a widget
    v-divider
    v-card-actions
      v-layout(justify-end)
        v-btn(@click="hideModal", depressed, flat) {{ $t('common.cancel') }}
</template>

<script>
import { MODALS, WIDGET_TYPES, SIDE_BARS_BY_WIDGET_TYPES } from '@/constants';
import { generateWidgetByType } from '@/helpers/entities';
import modalInnerMixin from '@/mixins/modal/inner';
import sideBarMixin from '@/mixins/side-bar/side-bar';

/**
 * Modal to create widget
 */
export default {
  name: MODALS.createWidget,
  mixins: [modalInnerMixin, sideBarMixin],
  data() {
    return {
      selectedWidget: '',
    };
  },
  computed: {
    widgetsTypes() {
      return [
        {
          title: this.$t('modals.widgetCreation.types.alarmList.title'),
          value: WIDGET_TYPES.alarmList,
          icon: 'view_list',
          category: 'general',
          // TODO: Add translation
          description: `Le widget Bac à alarmes permet de visualiser les alarmes présentes dans le SI.
            Il permet de visualiser les informations de l'alarme (composant concerné, statut de l'alarme, etc),
            ainsi que d'effectuer des actions sur ces alarmes.`,
        },
        {
          title: this.$t('modals.widgetCreation.types.context.title'),
          value: WIDGET_TYPES.context,
          icon: 'view_list',
          category: 'general',
          // TODO: Add translation
          description: `Le widget Explorateur de contexte permet de visualiser la liste des entités du SI.
            Il permet de visualiser les informations de l'entité, ainsi que d'effectuer des actions sur ces entités
            (ajout, suppression, édition)`,
        },
        {
          title: this.$t('modals.widgetCreation.types.weather.title'),
          value: WIDGET_TYPES.weather,
          icon: 'view_module',
          category: 'general',
          // TODO: Add translation
          description: `Le widget Météo de services permet de visualiser rapidement l'état des services.
            Celui-ci présente sous forme de grille un ensemble de services, avec un jeu de couleurs et d'icones
            permettant de rapidement visualiser une éventuelle source de problème.`,
        },
        {
          title: this.$t('modals.widgetCreation.types.statsHistogram.title'),
          value: WIDGET_TYPES.statsHistogram,
          icon: 'bar_chart',
          category: 'stat',
          // TODO: Add translation
          description: 'Le widget Histogramme permet de visualiser un ensemble de statistiques sous forme d\'histogramme.',
        },
        {
          title: this.$t('modals.widgetCreation.types.statsTable.title'),
          value: WIDGET_TYPES.statsTable,
          icon: 'table_chart',
          category: 'stat',
          // TODO: Add translation
          description: `Le widget Tableau de statistiques permet de visualiser sous forme de tableau un
            ensemble de statistiques du SI, ou d'un sous-ensemble du SI.`,
        },
        {
          title: this.$t('modals.widgetCreation.types.statsCalendar.title'),
          value: WIDGET_TYPES.statsCalendar,
          icon: 'calendar_today',
          category: 'stat',
          // TODO: Add translation
          description: `Le widget Calendrier permet de visualiser le nombre d'alarmes par période de temps
            (par heure, par jour, par mois), dans une vue calendrier.`,
        },
        {
          title: this.$t('modals.widgetCreation.types.statsCurves.title'),
          value: WIDGET_TYPES.statsCurves,
          icon: 'show_chart',
          category: 'stat',
          // TODO: Add translation
          description: 'Le widget Courbe permet de visualiser un ensemble de statistiques sous forme de courbes.',
        },
        {
          title: this.$t('modals.widgetCreation.types.statsNumber.title'),
          value: WIDGET_TYPES.statsNumber,
          icon: 'table_chart',
          category: 'stat',
          // TODO: Add translation
          description: `Le widget Compteur de statistique permet de visualiser la valeur d'une statistique,
            avec un jeu de couleur permettant de caractérisé la criticité de chacune des valeurs.`,
        },
        {
          title: this.$t('modals.widgetCreation.types.statsPareto.title'),
          value: WIDGET_TYPES.statsPareto,
          icon: 'multiline_chart',
          category: 'stat',
          // TODO: Add translation
          description: '',
        },
        {
          title: this.$t('modals.widgetCreation.types.text.title'),
          value: WIDGET_TYPES.text,
          icon: 'view_headline',
          category: 'general',
          // TODO: Add translation
          description: '',
        },
      ];
    },

    generalTypeWidgets() {
      return this.widgetsTypes.filter(widget => widget.category === 'general');
    },
    statTypeWidgets() {
      return this.widgetsTypes.filter(widget => widget.category === 'stat');
    },
  },
  methods: {
    selectWidgetType(type) {
      const widget = generateWidgetByType(type);

      this.showSideBar({
        name: SIDE_BARS_BY_WIDGET_TYPES[type],
        config: {
          widget,
          tabId: this.config.tabId,
          isNew: true,
        },
      });
      this.hideModal();
    },
    selectWidget(type) {
      this.selectedWidget = this.widgetsTypes.find(widgetType => widgetType.value === type);
    },
  },
};
</script>

<style lang="scss" scoped>
  .widgetType {
    cursor: pointer,
  }
</style>

