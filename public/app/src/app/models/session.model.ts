export class Session {
    name?: string;
    description?: string;
    date?: Date;
    startTime?: string;
    endTime?: string;

    constructor(data: any) {
        Object.assign(this, data);
    }
}
