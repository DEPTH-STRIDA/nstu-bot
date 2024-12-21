export interface Subject {
    ID: number;
    Name: string;
    CreatedAt: string;
    UpdatedAt: string;
    DeletedAt: null | string;
}

export interface LevelTraining {
    ID: number;
    Name: string;
    CreatedAt: string;
    UpdatedAt: string;
    DeletedAt: null | string;
}

export interface TutorExperience {
    ID: number;
    Name: string;
    value: number;
    CreatedAt: string;
    UpdatedAt: string;
    DeletedAt: null | string;
}

export interface WebSocketResponse {
    status: string;
    message: string;
}

export interface DataResponse extends WebSocketResponse {
    subjects: Subject[];
    level_trainings: LevelTraining[];
    level_experience: TutorExperience[];
}

export interface Teacher {
    ID: number;
    Name: string;
    Education: string;
    Experience?: TutorExperience;
    Subjects?: Subject[];
    Price: number;
    ClassFormat: string;
    ImgUrl: string;
    Services?: Service[];
}

export interface TeacherResponse extends WebSocketResponse {
    teachers: Teacher[];
    total: number;
    page: number;
    page_size: number;
}

export interface WebSocketRequest {
    action: string;
}

export interface TeacherRequest extends WebSocketRequest {
    filter: TeacherFilter;
}

export interface TeacherFilter {
    name: string;
    subjects: number[];
    level_training: number[];
    experience: number;
    price_from?: number;
    price_to?: number;
    page: number;
    page_size: number;
}

export interface Service {
    ID: number;
    Name: string;
} 