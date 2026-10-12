import { widgetAlarmsSocketMixin } from '@/mixins/widget/alarms-socket';

describe('widget-alarms-socket', () => {
  describe('methods', () => {
    test('Should join alarms socket room with alarm ids', () => {
      const addListener = vi.fn();
      const context = {
        $socket: {
          join: vi.fn(() => ({ addListener })),
        },
        alarmsSocketRoom: 'alarms/widget-id',
        updateAlarmInStore: vi.fn(),
      };

      widgetAlarmsSocketMixin.methods.joinToAlarmsSocketRoom.call(context, [
        { _id: 'alarm-id-1' },
        { _id: 'alarm-id-2' },
      ]);

      expect(context.$socket.join).toHaveBeenCalledWith('alarms/widget-id', [
        'alarm-id-1',
        'alarm-id-2',
      ]);
      expect(addListener).toHaveBeenCalledWith(context.updateAlarmInStore);
    });

    test('Should not join alarms socket room without alarms', () => {
      const context = {
        $socket: {
          join: vi.fn(),
        },
        alarmsSocketRoom: 'alarms/widget-id',
        updateAlarmInStore: vi.fn(),
      };

      widgetAlarmsSocketMixin.methods.joinToAlarmsSocketRoom.call(context, []);

      expect(context.$socket.join).not.toHaveBeenCalled();
    });
  });
});
