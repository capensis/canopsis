import { omit } from 'lodash';

import {
  generateViewTab,
  generateViewRow,
  generateWidgetByType,
} from '@/helpers/entities';

export default function copyViewTabs(tabs) {
  const widgetsIdsMap = {};

  const newTabs = tabs.map(tab => ({
    ...generateViewTab(),
    title: tab.title,
    rows: tab.rows.map(row => ({
      ...generateViewRow(),
      title: row.title,
      widgets: row.widgets.map((widget) => {
        const newWidget = generateWidgetByType(widget.type);

        widgetsIdsMap[widget._id] = newWidget._id; // Needed for userPreferences copy.

        return {
          ...newWidget,
          ...omit(widget, ['_id']),
        };
      }),
    })),
  }));

  return { tabs: newTabs, widgetsIdsMap };
}

export function copyTab(tab, title) {
  const widgetsIdsMap = {};

  const newTab = {
    ...generateViewTab(),
    title,
    rows: tab.rows.map(row => ({
      ...generateViewRow(),
      title: row.title,
      widgets: row.widgets.map((widget) => {
        const newWidget = generateWidgetByType(widget.type);

        widgetsIdsMap[widget._id] = newWidget._id; // Needed for userPreferences copy.

        return {
          ...newWidget,
          ...omit(widget, ['_id']),
        };
      }),
    })),
  };

  return { tab: newTab, widgetsIdsMap };
}
