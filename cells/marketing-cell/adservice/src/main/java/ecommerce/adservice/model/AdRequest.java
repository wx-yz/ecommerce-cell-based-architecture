package ecommerce.adservice.model;

import java.util.List;

public class AdRequest {
    private List<String> contextKeys;

    public AdRequest() {}

    public AdRequest(List<String> contextKeys) {
        this.contextKeys = contextKeys;
    }

    // Getters and Setters
    public List<String> getContextKeys() { return contextKeys; }
    public void setContextKeys(List<String> contextKeys) { this.contextKeys = contextKeys; }
}