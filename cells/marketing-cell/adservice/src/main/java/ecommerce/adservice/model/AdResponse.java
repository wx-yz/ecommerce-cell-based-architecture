package ecommerce.adservice.model;

import java.util.List;

public class AdResponse {
    private List<Ad> ads;

    public AdResponse() {}

    public AdResponse(List<Ad> ads) {
        this.ads = ads;
    }

    // Getters and Setters
    public List<Ad> getAds() { return ads; }
    public void setAds(List<Ad> ads) { this.ads = ads; }
}