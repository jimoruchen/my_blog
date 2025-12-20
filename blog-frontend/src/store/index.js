import { defineStore} from "pinia";
import axios from "axios";

export const useStore = defineStore('general',  {
    state: () => {
        return {
            user: {
                id: -1,
                username: '',
                email: '',
                role: '',
                avatar: null,
                registerTime: null
            },
            forum: {
                types: []
            }
        }
    }
})