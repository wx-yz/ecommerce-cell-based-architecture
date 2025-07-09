package ecommerce.adservice.model;

public class Ad {
    private String id;
    private String text;
    private String url;
    private String category;

    public Ad() {}

    public Ad(String id, String text, String url, String category) {
        this.id = id;
        this.text = text;
        this.url = url;
        this.category = category;
    }

    // Getters and Setters
    public String getId() { return id; }
    public void setId(String id) { this.id = id; }

    public String getText() { return text; }
    public void setText(String text) { this.text = text; }

    public String getUrl() { return url; }
    public void setUrl(String url) { this.url = url; }

    public String getCategory() { return category; }
    public void setCategory(String category) { this.category = category; }
}