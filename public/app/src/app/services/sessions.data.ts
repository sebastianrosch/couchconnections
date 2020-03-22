import { addDays } from 'date-fns';
import { Session } from '../models/session.model';

export const sessions: Session[] = [
    new Session({
        name: 'Aufrechterhaltung der Wirtschaft',
        description: 'Was können wir tun oder wie können wir den Kapitalismus temporär aussetzen?',
        date: new Date(2020, 2, 16),
        startTime: '18:00',
        endTime: '19:30'
    }),
    new Session({
        name: 'Das Leben nach Corona',
        description: 'Wie verändert sich die Welt?',
        date: new Date(2020, 2, 17),
        startTime: '20:00',
        endTime: '21:30'
    }),
    new Session({
        name: 'Remote Arbeiten',
        description: 'Wie läuft es bei euch? Kotzt ihr auch schon oder alles tutti?',
        date: new Date(2020, 2, 18),
        startTime: '19:00',
        endTime: '20:30'
    }),
    new Session({
        name: 'Ayurveda: Gesund trotz Hektik',
        description: '😌',
        date: new Date(2020, 2, 19),
        startTime: '21:00',
        endTime: '22:00'
    }),
    new Session({
        name: 'Entwicklung am FinTech Markt',
        description: 'Wer sind die Player und wie wird es das bestehende System erschüttern?',
        date: new Date(2020, 2, 20),
        startTime: '18:00',
        endTime: '20:00'
    }),
    new Session({
        name: 'Korruption',
        description: 'Das größte Gift der Gesellschaft. Wie können wir uns schützen?',
        date: new Date(2020, 2, 21),
        startTime: '20:00',
        endTime: '22:00'
    }),
    new Session({
        name: 'Home Schooling und Kinderbeschäftigung',
        description: 'Erfahrungsaustausch',
        date: new Date(2020, 2, 22),
        startTime: '14:00',
        endTime: '15:30'
    }),
    new Session({
        name: 'Archäologie',
        description: 'Lasst uns über die Relikte unserer Vergangenheit staunen!',
        date: new Date(2020, 2, 22),
        startTime: '18:00',
        endTime: '19:30'
    })
];
